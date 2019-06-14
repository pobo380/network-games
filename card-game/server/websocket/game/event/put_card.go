package event

import "github.com/pobo380/network-games/card-game/server/websocket/game/model"

type PutCard struct {
	PlayerId model.PlayerId
	Card     model.Card
}

func (*PutCard) GetType() string {
	return TypePutCard
}
