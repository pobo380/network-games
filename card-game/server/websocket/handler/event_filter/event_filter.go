package event_filter

import (
	"encoding/json"
	"github.com/pobo380/network-games/card-game/server/websocket/game/event"
	"github.com/pobo380/network-games/card-game/server/websocket/game/model"
)

type EventWithTypes []*EventWithType

type EventWithType struct {
	Type  string
	Event event.Event
}

func (ewt *EventWithType) UnmarshalJSON(b []byte) error {
	type alias EventWithType
	t := &struct {
		Event json.RawMessage
		*alias
	}{
		alias: (*alias)(ewt),
	}
	if err := json.Unmarshal(b, t); err != nil {
		return err
	}

	ev := event.NewFromType(t.Type)
	if err := json.Unmarshal(t.Event, ev); err != nil {
		return err
	}

	ewt.Event = ev
	return nil
}

func Filter(src event.Events, playerId string) (ret EventWithTypes) {
	for _, evt := range src {
		switch evt.GetType() {
		case event.TypeDrawCard:
			evt = filterDrawCard(evt, playerId)
		case event.TypeGameState:
			evt = filterGameState(evt, playerId)
		}

		ret = append(ret, &EventWithType{Type: evt.GetType(), Event: evt})
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
