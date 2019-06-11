package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

func OnDisconnect(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	keyCond := expression.Key("ConnectionId").Equal(expression.Value(request.RequestContext.ConnectionID))
	proj := expression.NamesList(expression.Name("PlayerId"))
	expr, err := expression.NewBuilder().WithProjection(proj).WithKeyCondition(keyCond).Build()
	if err != nil {
		return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 500}, err
	}

	q := &dynamodb.QueryInput{
		IndexName:                 aws.String(DynamoDbIndexConnectionId),
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		TableName:                 &DynamoDbTableConnections,
	}

	qr, err := dynamo.Query(q)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 500}, err
	}

	wrs := make([]*dynamodb.WriteRequest, 0, len(qr.Items))
	for _, item := range qr.Items {
		wrs = append(wrs, &dynamodb.WriteRequest{
			DeleteRequest: &dynamodb.DeleteRequest{
				Key: item,
			},
		})
	}

	bwi := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			DynamoDbTableConnections: wrs,
		},
	}

	_, err = dynamo.BatchWriteItem(bwi)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 200}, nil
}

func main() {
	lambda.Start(OnDisconnect)
}
