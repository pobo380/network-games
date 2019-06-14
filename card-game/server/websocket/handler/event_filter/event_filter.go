package event_filter

import (
	"github.com/pobo380/network-games/card-game/server/websocket/game/event"
	"github.com/pobo380/network-games/card-game/server/websocket/game/model"
)

func Filter(src event.Events, playerId string) (ret event.Events) {
	for _, evt := range src {
		var newEvt event.Event
		switch evt.GetType() {
		case event.TypeDrawCard:
			newEvt = filterDrawCard(evt, playerId)
		case event.TypeGameState:
			newEvt = filterGameState(evt, playerId)
		}

		if newEvt != nil {
			ret = append(ret, newEvt)
		}
	}

	return ret
}

func filterDrawCard(ev event.Event, playerId string) event.Event {
	dc := ev.(*event.DrawCard)

	if string(dc.PlayerId) == playerId {
		return ev
	}

	newEv := *dc
	newEv.Cards = model.Cards{}

	return &newEv
}

func filterGameState(ev event.Event, playerId string) event.Event {
	gs := ev.(*event.GameState)

	newEv := *gs
	for i := range newEv.State.Players {
		pl := &newEv.State.Players[i]
		if string(pl.Id) == playerId {
			continue
		}

		for j := range pl.Hand.Cards {
			pl.Hand.Cards[j] = model.InvalidCard
		}
	}

	return &newEv
}
