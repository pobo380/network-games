package response

import (
	"github.com/pobo380/network-games/card-game/server/websocket/handler/event_filter"
)

type GameEvent struct {
	Events event_filter.EventWithTypes
}
