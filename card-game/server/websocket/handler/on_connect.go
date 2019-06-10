package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pobo380/network-games/card-game/server/websocket/game/model"
	"github.com/pobo380/network-games/card-game/server/websocket/game/state"
)

func OnConnect(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	reqCtx := request.RequestContext

	in := &dynamodb.PutItemInput{
		TableName: &DynamoDbTableConnections,
		Item: map[string]*dynamodb.AttributeValue{
			"ConnectionId": {
				S: &reqCtx.ConnectionID,
			},
		},
	}
	_, err := dynamo.PutItem(in)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 500}, err
	}

	config := model.Config{
		InitialHandNum: 5,
	}
	players := []model.Player{
		{Id: "1"},
		{Id: "2"},
		{Id: "3"},
	}
	st := state.NewState(config, players)
	st.InitGame()

	av, err := dynamodbattribute.MarshalMap(st)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 500}, err
	}

	in2 := &dynamodb.PutItemInput{
		TableName: &DynamoDbTableGameStates,
		Item: map[string]*dynamodb.AttributeValue{
			"GameId": {
				S: &reqCtx.ConnectionID,
			},
			"State": {
				M: av,
			},
		},
	}
	_, err = dynamo.PutItem(in2)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 200}, nil
}

func main() {
	lambda.Start(OnConnect)
}
