package e2e

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pobo380/network-games/card-game/server/websocket/game/event"
	"github.com/pobo380/network-games/card-game/server/websocket/game/state"
	"github.com/pobo380/network-games/card-game/server/websocket/handler/request"
	"github.com/pobo380/network-games/card-game/server/websocket/handler/response"
	"github.com/pobo380/network-games/card-game/server/websocket/handler/table"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
)

func newWssConnection() *Client {
	playerId := uuid.NewV4().String()
	return newWssConnectionWithArgs(DefaultWssEndpoint, playerId)
}

func newWssConnectionWithArgs(url string, playerId string) *Client {
	h := http.Header{}
	h.Add("X-Pobo380-Network-Games-Player-Id", playerId)

	c, resp, err := websocket.DefaultDialer.Dial(url, h)
	if err != nil {
		b, _ := ioutil.ReadAll(resp.Body)
		panic(fmt.Sprintf("Dial failed : %s\n%+v\n%s", err, resp, string(b)))
	}

	return &Client{Con: c, PlayerId: playerId}
}

type Client struct {
	Con *websocket.Conn

	PlayerId string
	Room     *table.Room

	GameId string
	State  *state.State
}

func (c *Client) SendJoinRoom() {
	req := &request.Request{
		Action: "joinRoom",
		Body: &request.JoinRoomRequest{
			PlayerId: c.PlayerId,
		},
	}

	c.Con.WriteJSON(req)
}

func (c *Client) RecvRoomInfo() (*response.RoomInfo, error) {
	ri := &response.RoomInfo{}
	if err := c.Recv(ri); err != nil {
		return nil, err
	}

	if ri.Room != nil {
		c.Room = ri.Room
	}

	return ri, nil
}

func (c *Client) RecvGameStart() (*response.GameStart, error) {
	gs := &response.GameStart{}

	if err := c.Recv(gs); err != nil {
		return nil, err
	}

	c.GameId = gs.GameId

	return gs, nil
}

func (c *Client) RecvGameEvent() (*response.GameEvent, error) {
	ge := &response.GameEvent{}

	if err := c.Recv(ge); err != nil {
		return nil, err
	}

	for _, evt := range ge.Events {
		if evt.Type != "GameState" {
			continue
		}

		gs := evt.Event.(*event.GameState)
		c.State = gs.State
	}

	return ge, nil
}

func (c *Client) Recv(b interface{}) error {
	res := &response.Response{Body: b}
	return c.Con.ReadJSON(res)
}
