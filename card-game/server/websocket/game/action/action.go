package action

import (
	"github.com/pobo380/network-games/card-game/server/websocket/game/event"
	"github.com/pobo380/network-games/card-game/server/websocket/game/state"
)

const (
	TypePutCard = "PutCard"
	TypePass    = "Pass"
)

type Action interface {
	Do(st *state.State) (event.Events, *state.State)
}

func NewActionFromString(str string) Action {
	switch str {
	case TypePutCard:
		return &PutCard{}
	}
	return nil
}
