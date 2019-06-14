package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/private/protocol"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/pobo380/network-games/card-game/server/websocket/game/event"
	"github.com/pobo380/network-games/card-game/server/websocket/game/model"
	"github.com/pobo380/network-games/card-game/server/websocket/game/state"
	. "github.com/pobo380/network-games/card-game/server/websocket/handler"
	"github.com/pobo380/network-games/card-game/server/websocket/handler/event_filter"
	"github.com/pobo380/network-games/card-game/server/websocket/handler/request"
	"github.com/pobo380/network-games/card-game/server/websocket/handler/response"
	"github.com/pobo380/network-games/card-game/server/websocket/handler/table"
)

func JoinRoom(ctx context.Context, req events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	reqCtx := req.RequestContext
	commonRes := response.Responses{}

	// parse request
	jrr := &request.JoinRoomRequest{}
	err := request.Parse([]byte(req.Body), jrr)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: req.Body, StatusCode: 500}, err
	}

	// query open rooms
	keyCond := expression.Key("IsOpen").Equal(expression.Value("true"))
	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return events.APIGatewayProxyResponse{Body: req.Body, StatusCode: 500}, err
	}

	q := &dynamodb.QueryInput{
		IndexName:                 aws.String(DynamoDbIndexIsOpen),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Limit:                     aws.Int64(1),
		TableName:                 aws.String(DynamoDbTableRooms),
	}

	qr, err := Dynamo.Query(q)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: req.Body, StatusCode: 500}, err
	}

	// create or join room
	r := &table.Room{}
	if *qr.Count == 0 {
		// create Room
		rid := protocol.GetIdempotencyToken()
		r = &table.Room{
			RoomId:       rid,
			IsOpen:       "true",
			PlayerIds:    []string{jrr.PlayerId},
			MaxPlayerNum: MaxPlayerNumPerRoom,
		}
	} else {
		// join Room
		err = dynamodbattribute.UnmarshalMap(qr.Items[0], r)
		if err != nil {
			return events.APIGatewayProxyResponse{Body: req.Body, StatusCode: 500}, err
		}

		// add player and close when room is filled
		r.AddPlayer(jrr.PlayerId)
	}

	err = putItem(DynamoDbTableRooms, r)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: req.Body, StatusCode: 500}, err
	}

	// add RoomInfo response
	commonRes.Add(response.TypeRoomInfo, &response.RoomInfo{Room: r})

	// get pcs
	pcs, err := BatchGetPlayerConnections(r.PlayerIds)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: req.Body, StatusCode: 500}, err
	}

	// create apigw api manager
	gw, err := NewGwApi(reqCtx.DomainName, reqCtx.Stage)

	// send GameStart when room is filled
	if r.IsClosed() {
		st := newState(r.PlayerIds)
		raw, err := json.Marshal(st)
		if err != nil {
			return events.APIGatewayProxyResponse{Body: req.Body, StatusCode: 500}, err
		}

		// create Game
		g := &table.Game{
			GameId:    protocol.GetIdempotencyToken(),
			RawState:  string(raw),
			PlayerIds: r.PlayerIds,
		}

		err = putItem(DynamoDbTableGames, g)
		if err != nil {
			return events.APIGatewayProxyResponse{Body: req.Body, StatusCode: 500}, err
		}

		// send GameStart
		commonRes.Add(response.TypeGameStart, &response.GameStart{
			GameId:    g.GameId,
			PlayerIds: g.PlayerIds,
		})

		// create responses for each player
		resMap := make(map[string]response.Responses)
		for _, playerId := range g.PlayerIds {
			res := response.Responses{}
			res = append(res, commonRes...)

			// send GameEvent
			res.Add(response.TypeGameEvent, &response.GameEvent{
				Events: event_filter.Filter(event.Events{&event.GameState{State: st}}, playerId),
			})

			resMap[playerId] = res
		}

		err = SendResponsesMapToPlayers(gw, pcs, resMap)
		if err != nil {
			return events.APIGatewayProxyResponse{Body: req.Body, StatusCode: 500}, err
		}
	} else {
		err = SendResponsesToPlayers(gw, pcs, commonRes)
		if err != nil {
			return events.APIGatewayProxyResponse{Body: req.Body, StatusCode: 500}, err
		}
	}

	return events.APIGatewayProxyResponse{Body: req.Body, StatusCode: 200}, nil
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

	st := state.NewState(cfg, players)
	st.InitGame()

	return st
}

func putItem(table string, item interface{}) error {
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}

	_, err = Dynamo.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(table),
		Item:      av,
	})
	if err != nil {
		return err
	}

	return nil
}

func main() {
	lambda.Start(JoinRoom)
}
