package main

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	gwApi "github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pobo380/network-games/card-game/server/websocket/handler/response"
	"github.com/pobo380/network-games/card-game/server/websocket/table"
)

func batchGetPlayerConnections(playerIds []string) ([]*table.PlayerConnection, error) {
	avs := make([]map[string]*dynamodb.AttributeValue, 0, len(playerIds))
	for _, playerId := range playerIds {
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
		return nil, err
	}

	responses := items.Responses[DynamoDbTableConnections]
	pcs := make([]*table.PlayerConnection, 0, len(responses))
	for _, item := range responses {
		pc := &table.PlayerConnection{}
		err := dynamodbattribute.UnmarshalMap(item, pc)
		if err != nil {
			return nil, err
		}
		pcs = append(pcs, pc)
	}

	return pcs, nil
}

func sendResponsesToPlayers(gw *gwApi.ApiGatewayManagementApi, pcs []*table.PlayerConnection, res response.Responses) error {
	for _, r := range res {
		raw, err := json.Marshal(r)
		if err != nil {
			return err
		}

		for _, pc := range pcs {
			data := &gwApi.PostToConnectionInput{
				ConnectionId: aws.String(pc.ConnectionId),
				Data:         raw,
			}

			_, err = gw.PostToConnection(data)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
