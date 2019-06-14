package response

import "github.com/pobo380/network-games/card-game/server/websocket/handler/table"

type RoomInfo struct {
	Room *table.Room
}
