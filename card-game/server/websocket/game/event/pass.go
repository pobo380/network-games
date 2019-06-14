package event

import "github.com/pobo380/network-games/card-game/server/websocket/game/model"

type Pass struct {
	PlayerId model.PlayerId
}

func (*Pass) GetType() string {
	return TypePass
}
