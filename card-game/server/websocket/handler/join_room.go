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
	"github.com/pobo380/network-games/card-game/server/websocket/room"
)

type JoinRoomRequest struct {
	PlayerId string
}

func JoinRoom(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
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
	if *qr.Count == 0 {
		// create Room
		rid := protocol.GetIdempotencyToken()
		r := &room.Room{
			RoomId:       rid,
			IsOpen:       "true",
			PlayerIds:    []string{payload.PlayerId},
			MaxPlayerNum: MaxPlayerNumPerRoom,
		}

		err = putItemRoom(r)
		if err != nil {
			return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 500}, err
		}
	} else {
		// join Room
		r := &room.Room{}
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
	}

	return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 200}, nil
}

func putItemRoom(r *room.Room) error {
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

func main() {
	lambda.Start(JoinRoom)
}
