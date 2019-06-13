package request

import "encoding/json"

type GameActionRequest struct {
	Type       string
	GameId     string
	GameAction *json.RawMessage
}
