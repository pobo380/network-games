package event

import "github.com/pobo380/network-games/card-game/server/websocket/game/model"

type FinishType int

const (
	FinishTypeWin FinishType = iota
	FinishTypeDraw
)

type FinishGame struct {
	WinnerId   model.PlayerId
	FinishType FinishType
}

func (*FinishGame) GetType() string {
	return TypeFinishGame
}
