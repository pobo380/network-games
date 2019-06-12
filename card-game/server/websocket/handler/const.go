package main

const (
	EnvDynamoDbTableConnections = "DYNAMODB_TABLE_CONNECTIONS"
	EnvDynamoDbTableRooms       = "DYNAMODB_TABLE_ROOMS"
	EnvDynamoDbTableGameStates  = "DYNAMODB_TABLE_GAMES"

	DynamoDbIndexConnectionId = "ConnectionId"
	DynamoDbIndexIsOpen       = "IsOpen"

	CustomHttpHeaderPrefix = "x-pobo380-network-games"
	CustomHeaderPlayerId   = CustomHttpHeaderPrefix + "-player-id"

	MaxPlayerNumPerRoom = 4
)
