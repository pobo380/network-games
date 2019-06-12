package action

import (
	"github.com/pobo380/network-games/card-game/server/websocket/game/event"
	"github.com/pobo380/network-games/card-game/server/websocket/game/model"
	"github.com/pobo380/network-games/card-game/server/websocket/game/state"
)

type Pass struct {
	PlayerId model.PlayerId
}

func (p *Pass) Do(st *state.State) (events event.Events, retState *state.State) {
	retState = st
	errFn := func(reason string) (event.Events, *state.State) {
		events.Add(&event.NotAcceptable{
			ActionType: string(TypePass),
			Reason:     reason,
		})
		return events, retState
	}

	player := st.FindPlayerById(p.PlayerId)
	if player == nil {
		return errFn("PlayerId is not found.")
	}

	if st.PlayOrder.CurrentPlayerId() != p.PlayerId {
		return errFn("PlayerId is not turn owner.")
	}

	// Pass
	events.Add(&event.Pass{
		PlayerId: p.PlayerId,
	})

	// DrawCard
	if st.Deck.Empty() {
		if st.Discards.Empty() {
			events.Add(&event.FinishGame{
				FinishType: event.FinishTypeDraw,
			})
			return
		}

		st.ResetDeck()
	}

	c, _ := st.Deck.Draw()
	player.Hand.Add(c)

	events.Add(&event.DrawCard{
		PlayerId: p.PlayerId,
		Cards:    model.Cards{c},
		CardsNum: 1,
	})

	// TurnChange
	nextPlayerId := st.PlayOrder.StepToNextPlayer()
	events.Add(&event.TurnChange{PlayerId: nextPlayerId})

	return
}
