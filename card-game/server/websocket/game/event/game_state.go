package event

import "github.com/pobo380/network-games/card-game/server/websocket/game/state"

type GameState struct {
	State *state.State
}

func (*GameState) GetType() string {
	return TypeGameState
}
