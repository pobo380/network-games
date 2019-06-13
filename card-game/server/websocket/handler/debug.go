package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

func Debug(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	in := &dynamodb.ListTablesInput{}
	out, err := dynamo.ListTables(in)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: req.Body, StatusCode: 500}, err
	}

	for _, tbl := range out.TableNames {
		descIn := &dynamodb.DescribeTableInput{
			TableName: tbl,
		}

		descOut, err := dynamo.DescribeTable(descIn)
		if err != nil {
			fmt.Printf("Skip table : %s\n", *tbl)
			continue
		}

		fmt.Printf("Delete table : %s\n", *tbl)

		nb := make([]expression.NameBuilder, 0)
		for _, ks := range descOut.Table.KeySchema {
			nb = append(nb, expression.Name(*ks.AttributeName))
		}

		proj := expression.NamesList(nb[0], nb[1:]...)
		expr, err := expression.NewBuilder().WithProjection(proj).Build()
		if err != nil {
			return events.APIGatewayProxyResponse{Body: req.Body, StatusCode: 500}, err
		}

		scIn := &dynamodb.ScanInput{
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			ProjectionExpression:      expr.Projection(),
			TableName:                 tbl,
		}

		scOut, err := dynamo.Scan(scIn)
		if err != nil {
			return events.APIGatewayProxyResponse{Body: req.Body, StatusCode: 500}, err
		}
		fmt.Printf("  Scan success\n")

		for _, item := range scOut.Items {
			fmt.Printf("    Delete item : %v\n", item)
			delIn := &dynamodb.DeleteItemInput{
				Key:       item,
				TableName: tbl,
			}

			_, err := dynamo.DeleteItem(delIn)
			if err != nil {
				return events.APIGatewayProxyResponse{Body: req.Body, StatusCode: 500}, err
			}
		}
	}

	return events.APIGatewayProxyResponse{Body: req.Body, StatusCode: 200}, nil
}

func main() {
	lambda.Start(Debug)
}
