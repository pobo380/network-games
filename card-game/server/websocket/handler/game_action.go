package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pobo380/network-games/card-game/server/websocket/game/action"
	"github.com/pobo380/network-games/card-game/server/websocket/game/state"
	"github.com/pobo380/network-games/card-game/server/websocket/handler/response"
	"github.com/pobo380/network-games/card-game/server/websocket/table"
)

type GameActionRequest struct {
	Type       string
	GameId     string
	GameAction *json.RawMessage
}

func GameAction(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	reqCtx := request.RequestContext
	res := response.Responses{}

	gar := &GameActionRequest{}
	err := json.Unmarshal([]byte(request.Body), gar)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 500}, err
	}

	act := action.NewActionFromType(action.Type(gar.Type))
	err = json.Unmarshal(*gar.GameAction, act)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 500}, err
	}

	// get Game
	game, err := getItemGame(gar.GameId)

	st := &state.State{}
	err = json.Unmarshal([]byte(game.RawState), st)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 500}, err
	}

	// do action
	evs, st := act.Do(st)

	// save Game
	rawSt, err := json.Marshal(st)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 500}, err
	}
	game.RawState = string(rawSt)

	err = putItemGame(game)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 500}, err
	}

	// send Events
	res.Add(response.TypeGameEvent, &response.GameEvent{Events: evs})

	// create apigw api manager
	gw, err := NewGwApi(reqCtx.DomainName, reqCtx.Stage)

	// send Responses
	pcs, err := batchGetPlayerConnections(game.PlayerIds)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 500}, err
	}

	err = sendResponsesToPlayers(gw, pcs, res)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 200}, nil
}

func getItemGame(gameId string) (*table.Game, error) {
	in := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"GameId": {
				S: aws.String(gameId),
			},
		},
		TableName: aws.String(DynamoDbTableGames),
	}
	out, err := dynamo.GetItem(in)
	if err != nil {
		return nil, err
	}

	game := &table.Game{}
	err = dynamodbattribute.UnmarshalMap(out.Item, game)
	if err != nil {
		return nil, err
	}

	return game, nil
}

func putItemGame(g *table.Game) error {
	av, err := dynamodbattribute.MarshalMap(g)
	if err != nil {
		return err
	}

	in := &dynamodb.PutItemInput{
		TableName: aws.String(DynamoDbTableGames),
		Item:      av,
	}

	_, err = dynamo.PutItem(in)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	lambda.Start(GameAction)
}
