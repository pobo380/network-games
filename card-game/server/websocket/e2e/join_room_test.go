package e2e

import (
	"github.com/pobo380/network-games/card-game/server/websocket/handler/request"
	"github.com/pobo380/network-games/card-game/server/websocket/handler/response"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_JoinRoom(t *testing.T) {
	con, playerId := newWssConnection()
	defer con.Close()

	req := &request.Request{
		Action: "joinRoom",
		Body: &request.JoinRoomRequest{
			PlayerId: playerId,
		},
	}

	con.WriteJSON(req)

	ri := &response.RoomInfo{}
	res := &response.Response{Body: ri}
	err := con.ReadJSON(res)

	assert.NoError(t, err)
	assert.NotNil(t, ri.Room)
}
