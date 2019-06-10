package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func OnDisconnect(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	reqCtx := request.RequestContext

	in := &dynamodb.DeleteItemInput{
		TableName: &DynamoDbTableConnections,
		Key: map[string]*dynamodb.AttributeValue{
			"ConnectionId": {
				S: &reqCtx.ConnectionID,
			},
		},
	}
	_, err := dynamo.DeleteItem(in)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 200}, nil
}

func main() {
	lambda.Start(OnDisconnect)
}
