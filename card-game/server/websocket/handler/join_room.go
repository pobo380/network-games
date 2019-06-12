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
	"github.com/pobo380/network-games/card-game/server/websocket/handler/response"
	"github.com/pobo380/network-games/card-game/server/websocket/table"
)

type JoinRoomRequest struct {
	PlayerId string
}

func JoinRoom(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	reqCtx := request.RequestContext

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

	var res *table.Room

	// create or join room
	if *qr.Count == 0 {
		// create Room
		rid := protocol.GetIdempotencyToken()
		r := &table.Room{
			RoomId:       rid,
			IsOpen:       "true",
			PlayerIds:    []string{payload.PlayerId},
			MaxPlayerNum: MaxPlayerNumPerRoom,
		}

		err = putItemRoom(r)
		if err != nil {
			return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 500}, err
		}

		res = r
	} else {
		// join Room
		r := &table.Room{}
		err = dynamodbattribute.UnmarshalMap(qr.Items[0], r)
		if err != nil {
			return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 500}, err
		}

		// add player and close when room is filled
		r.AddPlayer(payload.PlayerId)

		err = putItemRoom(r)
		if err != nil {
			return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 500}, err
		}

		res = r
	}

	// create apigw api manager
	gw, err := NewGwApi(reqCtx.DomainName, reqCtx.Stage)

	// send RoomInfo
	sendRoomInfoToPlayers(gw, &response.RoomInfo{Room: res})

	return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 200}, nil
}

func putItemRoom(r *table.Room) error {
	item, err := dynamodbattribute.MarshalMap(r)
	if err != nil {
		return err
	}

	_, err = dynamo.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(DynamoDbTableRooms),
		Item:      item,
	})
	if err != nil {
		return err
	}

	return nil
}

func sendRoomInfoToPlayers(gw *gwApi.ApiGatewayManagementApi, info *response.RoomInfo) error {
	avs := make([]map[string]*dynamodb.AttributeValue, 0, len(info.Room.PlayerIds))
	for _, playerId := range info.Room.PlayerIds {
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
		return err
	}

	res := &response.Response{
		Type: response.TypeRoomInfo,
		Body: info,
	}

	raw, err := json.Marshal(res)
	if err != nil {
		return err
	}

	for _, item := range items.Responses[DynamoDbTableConnections] {
		pc := &table.PlayerConnection{}
		err := dynamodbattribute.UnmarshalMap(item, pc)
		if err != nil {
			return err
		}

		data := &gwApi.PostToConnectionInput{
			ConnectionId: aws.String(pc.ConnectionId),
			Data:         raw,
		}

		_, err = gw.PostToConnection(data)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	lambda.Start(JoinRoom)
}
