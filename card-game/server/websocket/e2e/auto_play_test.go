package e2e

import (
	"github.com/pobo380/network-games/card-game/server/websocket/handler/request"
	"github.com/pobo380/network-games/card-game/server/websocket/handler/response"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_AutoPlay(t *testing.T) {
	cA := newWssConnection()
	defer cA.Con.Close()

	cB := newWssConnection()
	defer cB.Con.Close()

	cC := newWssConnection()
	defer cC.Con.Close()

	cD := newWssConnection()
	defer cD.Con.Close()

	sendJoinRoomReq := func(c *Client) {
		req := &request.Request{
			Action: "joinRoom",
			Body: &request.JoinRoomRequest{
				PlayerId: c.PlayerId,
			},
		}

		c.Con.WriteJSON(req)

		time.Sleep(500 * time.Millisecond)
	}

	recvRoomInfo := func(c *Client) string {
		ri := &response.RoomInfo{}
		res := &response.Response{Body: ri}
		err := c.Con.ReadJSON(res)

		assert.NoError(t, err)
		assert.NotNil(t, ri.Room)
		assert.NotEmpty(t, ri.Room.RoomId)

		return ri.Room.RoomId
	}

	// join A
	sendJoinRoomReq(cA)
	recvRoomInfo(cA)

	// join B
	sendJoinRoomReq(cB)
	recvRoomInfo(cA)
	recvRoomInfo(cB)

	// join C
	sendJoinRoomReq(cC)
	recvRoomInfo(cA)
	recvRoomInfo(cB)
	recvRoomInfo(cC)

	// join D
	sendJoinRoomReq(cD)
	recvRoomInfo(cA)
	recvRoomInfo(cB)
	recvRoomInfo(cC)
	recvRoomInfo(cD)

	recvStartGame := func(c *Client) string {
		gs := &response.GameStart{}
		res := &response.Response{Body: gs}
		err := c.Con.ReadJSON(res)

		assert.NoError(t, err)
		assert.NotEmpty(t, gs.GameId)

		return gs.GameId
	}

	func() {
		recvStartGame(cA)
		recvStartGame(cB)
		recvStartGame(cC)
		recvStartGame(cD)
	}()

	recvGameEvent := func(c *Client) {
		ge := &response.GameEvent{}
		res := &response.Response{Body: ge}
		err := c.Con.ReadJSON(res)

		assert.NoError(t, err)
		assert.Equal(t, ge.Events[0].Type, "GameState")
	}

	func() {
		recvGameEvent(cA)
		recvGameEvent(cB)
		recvGameEvent(cC)
		recvGameEvent(cD)
	}()
}
