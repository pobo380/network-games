package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	gwApi "github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

func SendMessage(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	reqCtx := request.RequestContext

	// scan ConnectionId
	proj := expression.NamesList(expression.Name("ConnectionId"))
	expr, err := expression.NewBuilder().WithProjection(proj).Build()
	if err != nil {
		return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 500}, err
	}

	in := &dynamodb.ScanInput{
		TableName:                &DynamoDbTableConnections,
		ExpressionAttributeNames: expr.Names(),
		ProjectionExpression:     expr.Projection(),
	}
	out, err := dynamo.Scan(in)

	// create apigw api manager
	gw, err := NewGwApi(reqCtx.DomainName, reqCtx.Stage)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 500}, err
	}

	// send to all ConnectionId
	for _, item := range out.Items {
		connectionId := item["ConnectionId"].S

		// post to connection
		data := &gwApi.PostToConnectionInput{
			ConnectionId: connectionId,
			Data:         []byte(request.Body),
		}
		gw.PostToConnection(data)
	}

	return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 200}, nil
}

func main() {
	lambda.Start(SendMessage)
}
