package event

import "github.com/pobo380/network-games/card-game/server/websocket/game/model"

type TurnChange struct {
	PlayerId model.PlayerId
}

func (*TurnChange) GetType() string {
	return TypeTurnChange
}
