package action

import (
	"github.com/pobo380/network-games/card-game/server/websocket/game/event"
	"github.com/pobo380/network-games/card-game/server/websocket/game/model"
	"github.com/pobo380/network-games/card-game/server/websocket/game/state"
)

type PutCard struct {
	PlayerId model.PlayerId
	Card     model.Card
}

func (pc *PutCard) Do(st *state.State) (events event.Events, retState *state.State) {
	retState = st
	errFn := func(reason string) (event.Events, *state.State) {
		events.Add(&event.NotAcceptable{
			ActionType: TypePutCard,
			Reason:     reason,
		})
		return events, retState
	}

	player := st.FindPlayerById(pc.PlayerId)
	if player == nil {
		return errFn("PlayerId is not found.")
	}

	if st.PlayOrder.CurrentPlayerId() != pc.PlayerId {
		return errFn("PlayerId is not turn owner.")
	}

	top, exists := st.Upcards.Top()
	if !exists {
		return errFn("No upcards exists.")
	}

	if !top.SameNumber(pc.Card) && !top.SameSuit(pc.Card) {
		return errFn("PutCard has different number and different suit.")
	}

	if !player.Hand.Contain(pc.Card) {
		return errFn("Player has not PutCard.")
	}

	// PutCard
	player.Hand.Discard(pc.Card)
	st.Upcards.Put(pc.Card)
	st.Discards.Put(top)

	events.Add(&event.PutCard{PlayerId: pc.PlayerId, Card: pc.Card})

	// FinishGame
	if player.Hand.Empty() {
		events.Add(&event.FinishGame{
			WinnerId:   pc.PlayerId,
			FinishType: event.FinishTypeWin,
		})
		return
	}

	// TurnChange
	nextPlayerId := st.PlayOrder.StepToNextPlayer()
	events.Add(&event.TurnChange{PlayerId: nextPlayerId})

	return
}
