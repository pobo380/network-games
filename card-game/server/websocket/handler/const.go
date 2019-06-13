package main

const (
	EnvDynamoDbTableConnections = "DYNAMODB_TABLE_CONNECTIONS"
	EnvDynamoDbTableRooms       = "DYNAMODB_TABLE_ROOMS"
	EnvDynamoDbTableGameStates  = "DYNAMODB_TABLE_GAMES"

	DynamoDbIndexConnectionId = "ConnectionId"
	DynamoDbIndexIsOpen       = "IsOpen"

	CustomHttpHeaderPrefix = "X-Pobo380-Network-Games"
	CustomHeaderPlayerId   = CustomHttpHeaderPrefix + "-Player-Id"

	MaxPlayerNumPerRoom = 4
)
