package handler

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	gwApi "github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var (
	Dynamo = NewDynamo()

	DynamoDbTableConnections = os.Getenv(EnvDynamoDbTableConnections)
	DynamoDbTableRooms       = os.Getenv(EnvDynamoDbTableRooms)
	DynamoDbTableGames       = os.Getenv(EnvDynamoDbTableGameStates)
)

func NewDynamo() *dynamodb.DynamoDB {
	sess, _ := session.NewSession()
	return dynamodb.New(sess)
}

func NewGwApi(domainName string, stage string) (*gwApi.ApiGatewayManagementApi, error) {
	endpoint := fmt.Sprintf("%s/%s", domainName, stage)
	sess, err := session.NewSession(&aws.Config{
		Endpoint: &endpoint,
	})
	if err != nil {
		return nil, err
	}

	return gwApi.New(sess), nil
}
