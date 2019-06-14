package event

import "github.com/pobo380/network-games/card-game/server/websocket/game/model"

type DrawCard struct {
	PlayerId model.PlayerId
	Cards    model.Cards
	CardsNum int
}

func (*DrawCard) GetType() string {
	return TypeDrawCard
}
