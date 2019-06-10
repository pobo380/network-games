package state

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"

	"github.com/pobo380/network-games/card-game/server/websocket/game/model"
)

func TestState_Deal(t *testing.T) {
	type fields struct {
		Rand     *rand.Rand
		Config   model.Config
		Players  []model.Player
		Deck     model.Deck
		Discards model.Discards
		Upcards  model.Upcards
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "正常系",
			fields: fields{
				Rand: rand.New(rand.NewSource(100)),
				Config: model.Config{
					InitialHandNum: 3,
				},
				Players: []model.Player{
					{
						Id: "1",
					},
					{
						Id: "2",
					},
					{
						Id: "3",
					},
					{
						Id: "4",
					},
				},
				Deck: model.Deck{
					Cards: model.Cards{
						{Suit: model.SuitHeart, Number: 1},
						{Suit: model.SuitHeart, Number: 2},
						{Suit: model.SuitHeart, Number: 3},
						{Suit: model.SuitHeart, Number: 4},
						{Suit: model.SuitHeart, Number: 5},
						{Suit: model.SuitHeart, Number: 6},
						{Suit: model.SuitHeart, Number: 7},
						{Suit: model.SuitHeart, Number: 8},
						{Suit: model.SuitHeart, Number: 9},
						{Suit: model.SuitDiamond, Number: 1},
						{Suit: model.SuitDiamond, Number: 2},
						{Suit: model.SuitDiamond, Number: 3},
						{Suit: model.SuitDiamond, Number: 4},
						{Suit: model.SuitDiamond, Number: 5},
						{Suit: model.SuitDiamond, Number: 6},
						{Suit: model.SuitDiamond, Number: 7},
						{Suit: model.SuitDiamond, Number: 8},
						{Suit: model.SuitDiamond, Number: 9},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &State{
				rand:     tt.fields.Rand,
				Config:   tt.fields.Config,
				Players:  tt.fields.Players,
				Deck:     tt.fields.Deck,
				Discards: tt.fields.Discards,
				Upcards:  tt.fields.Upcards,
			}
			st.Deal()

			handExpects := []model.Cards{
				{
					{Suit: model.SuitHeart, Number: 1},
					{Suit: model.SuitHeart, Number: 2},
					{Suit: model.SuitHeart, Number: 3},
				},
				{
					{Suit: model.SuitHeart, Number: 4},
					{Suit: model.SuitHeart, Number: 5},
					{Suit: model.SuitHeart, Number: 6},
				},
				{
					{Suit: model.SuitHeart, Number: 7},
					{Suit: model.SuitHeart, Number: 8},
					{Suit: model.SuitHeart, Number: 9},
				},
				{
					{Suit: model.SuitDiamond, Number: 1},
					{Suit: model.SuitDiamond, Number: 2},
					{Suit: model.SuitDiamond, Number: 3},
				},
			}

			deckExpects := model.Cards{
				{Suit: model.SuitDiamond, Number: 4},
				{Suit: model.SuitDiamond, Number: 5},
				{Suit: model.SuitDiamond, Number: 6},
				{Suit: model.SuitDiamond, Number: 7},
				{Suit: model.SuitDiamond, Number: 8},
				{Suit: model.SuitDiamond, Number: 9},
			}

			assert.Equal(t, deckExpects, st.Deck.Cards)

			for i, expect := range handExpects {
				assert.Equal(t, expect, st.Players[i].Hand.Cards)
			}
		})
	}
}
