package action

import (
	"github.com/pobo380/network-games/card-game/server/websocket/game/event"
	"github.com/pobo380/network-games/card-game/server/websocket/game/state"
)

type Type string

const (
	TypePutCard Type = "PutCard"
	TypePass    Type = "Pass"
)

type Action interface {
	Do(st *state.State) (event.Events, *state.State)
}

func NewActionFromType(t Type) Action {
	switch t {
	case TypePutCard:
		return &PutCard{}
	case TypePass:
		return &Pass{}
	}
	return nil
}
