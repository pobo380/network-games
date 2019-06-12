package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/private/protocol"
	gwApi "github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/pobo380/network-games/card-game/server/websocket/game/model"
	"github.com/pobo380/network-games/card-game/server/websocket/game/state"
	"github.com/pobo380/network-games/card-game/server/websocket/handler/response"
	"github.com/pobo380/network-games/card-game/server/websocket/table"
)

type JoinRoomRequest struct {
	PlayerId string
}

func JoinRoom(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	reqCtx := request.RequestContext
	res := response.Responses{}

	// parse request
	payload := &JoinRoomRequest{}
	err := json.Unmarshal([]byte(request.Body), payload)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 500}, err
	}

	// query open rooms
	keyCond := expression.Key("IsOpen").Equal(expression.Value("true"))
	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 500}, err
	}

	q := &dynamodb.QueryInput{
		IndexName:                 aws.String(DynamoDbIndexIsOpen),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Limit:                     aws.Int64(1),
		TableName:                 aws.String(DynamoDbTableRooms),
	}

	qr, err := dynamo.Query(q)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 500}, err
	}

	// create or join room
	r := &table.Room{}
	if *qr.Count == 0 {
		// create Room
		rid := protocol.GetIdempotencyToken()
		r = &table.Room{
			RoomId:       rid,
			IsOpen:       "true",
			PlayerIds:    []string{payload.PlayerId},
			MaxPlayerNum: MaxPlayerNumPerRoom,
		}
	} else {
		// join Room
		err = dynamodbattribute.UnmarshalMap(qr.Items[0], r)
		if err != nil {
			return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 500}, err
		}

		// add player and close when room is filled
		r.AddPlayer(payload.PlayerId)
	}

	err = putItem(DynamoDbTableRooms, r)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 500}, err
	}

	// add RoomInfo response
	res.Add(response.TypeRoomInfo, &response.RoomInfo{Room: r})

	// send GameStart when room is filled
	if r.IsClosed() {
		st := newState(r.PlayerIds)
		raw, err := json.Marshal(st)
		if err != nil {
			return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 500}, err
		}

		// create Game
		g := &table.Game{
			GameId:    protocol.GetIdempotencyToken(),
			RawState:  string(raw),
			PlayerIds: r.PlayerIds,
		}

		err = putItem(DynamoDbTableGames, g)
		if err != nil {
			return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 500}, err
		}

		// send GameStart
		res.Add(response.TypeGameStart, &response.GameStart{
			GameId:    g.GameId,
			PlayerIds: g.PlayerIds,
		})
	}

	// create apigw api manager
	gw, err := NewGwApi(reqCtx.DomainName, reqCtx.Stage)

	// send Responses
	pcs, err := batchGetPlayerConnections(r.PlayerIds)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 500}, err
	}

	err = sendResponsesToPlayers(gw, pcs, res)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 200}, nil
}

func newState(playerIds []string) *state.State {
	cfg := model.Config{
		InitialHandNum: 5,
	}

	players := make([]model.Player, 0, len(playerIds))
	for _, pid := range playerIds {
		players = append(players, model.Player{
			Id:   model.PlayerId(pid),
			Hand: model.Hand{},
		})
	}

	return state.NewState(cfg, players)
}

func putItem(table string, item interface{}) error {
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}

	_, err = dynamo.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(table),
		Item:      av,
	})
	if err != nil {
		return err
	}

	return nil
}

func batchGetPlayerConnections(playerIds []string) ([]*table.PlayerConnection, error) {
	avs := make([]map[string]*dynamodb.AttributeValue, 0, len(playerIds))
	for _, playerId := range playerIds {
		avs = append(avs, map[string]*dynamodb.AttributeValue{
			"PlayerId": {
				S: aws.String(playerId),
			},
		})
	}

	bgi := &dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			DynamoDbTableConnections: {
				Keys: avs,
			},
		},
	}

	items, err := dynamo.BatchGetItem(bgi)
	if err != nil {
		return nil, err
	}

	responses := items.Responses[DynamoDbTableConnections]
	pcs := make([]*table.PlayerConnection, 0, len(responses))
	for _, item := range responses {
		pc := &table.PlayerConnection{}
		err := dynamodbattribute.UnmarshalMap(item, pc)
		if err != nil {
			return nil, err
		}
		pcs = append(pcs, pc)
	}

	return pcs, nil
}

func sendResponsesToPlayers(gw *gwApi.ApiGatewayManagementApi, pcs []*table.PlayerConnection, res response.Responses) error {
	for _, r := range res {
		raw, err := json.Marshal(r)
		if err != nil {
			return err
		}

		for _, pc := range pcs {
			data := &gwApi.PostToConnectionInput{
				ConnectionId: aws.String(pc.ConnectionId),
				Data:         raw,
			}

			_, err = gw.PostToConnection(data)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func main() {
	lambda.Start(JoinRoom)
}
