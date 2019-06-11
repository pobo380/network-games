package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pobo380/network-games/card-game/server/websocket/connection"
)

func OnConnect(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	reqCtx := request.RequestContext
	playerId := request.Headers[CustomHeaderPlayerId]

	pc := &connection.PlayerConnection{
		PlayerId:     playerId,
		ConnectionId: reqCtx.ConnectionID,
	}

	item, err := dynamodbattribute.MarshalMap(pc)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 500}, err
	}

	in := &dynamodb.PutItemInput{
		TableName: &DynamoDbTableConnections,
		Item:      item,
	}

	_, err = dynamo.PutItem(in)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 200}, nil
}

func main() {
	lambda.Start(OnConnect)
}
