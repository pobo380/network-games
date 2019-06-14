package e2e

import (
	"fmt"
	"github.com/pobo380/network-games/card-game/server/websocket/game/action"
	"github.com/pobo380/network-games/card-game/server/websocket/game/event"
	"github.com/pobo380/network-games/card-game/server/websocket/game/model"
	"github.com/pobo380/network-games/card-game/server/websocket/handler/response"
	"github.com/stretchr/testify/assert"
	"testing"
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

	clients := []*Client{cA, cB, cC, cD}

	// join A
	cA.SendJoinRoom()
	cA.RecvRoomInfo()

	// join B
	cB.SendJoinRoom()
	cA.RecvRoomInfo()
	cB.RecvRoomInfo()

	// join C
	cC.SendJoinRoom()
	cA.RecvRoomInfo()
	cB.RecvRoomInfo()
	cC.RecvRoomInfo()

	// join D
	cD.SendJoinRoom()
	for _, c := range clients {
		c.RecvRoomInfo()
	}

	// recv GameStart
	for _, c := range clients {
		c.RecvGameStart()
	}

	// recv GameEvent
	for _, c := range clients {
		c.RecvGameEvent()
	}

	assert.NotNil(t, cA.State)
	assert.NotNil(t, cB.State)
	assert.NotNil(t, cC.State)
	assert.NotNil(t, cD.State)

	findClientAndPlayer := func(playerId model.PlayerId) (*Client, *model.Player) {
		for _, c := range clients {
			if c.PlayerId == string(playerId) {
				return c, c.State.FindPlayerById(playerId)
			}
		}
		return nil, nil
	}

	gameFinished := func(ge *response.GameEvent) bool {
		for _, ewt := range ge.Events {
			evt := ewt.Event

			switch evt.(type) {
			case *event.FinishGame:
				return true
			}
		}
		return false
	}

	turn := 0
	for {
		playerId := cA.State.PlayOrder.CurrentPlayerId()
		cl, player := findClientAndPlayer(playerId)

		uc, _ := cA.State.Upcards.Top()

		fmt.Printf("Turn Change\n")
		fmt.Printf("  PlayerId : %s\n", playerId)
		fmt.Printf("  Hand     : %+v\n", player.Hand.Cards)
		fmt.Printf("  Upcard   : %+v\n", uc)
		fmt.Println("")

		found := false
		for _, hc := range player.Hand.Cards {
			if uc.SameNumber(hc) || uc.SameSuit(hc) {
				req := &action.PutCard{
					PlayerId: playerId,
					Card:     hc,
				}
				cl.SendGameAction(action.TypePutCard, req)

				fmt.Printf("Put Card\n")
				fmt.Printf("  Card : %+v\n", hc)
				fmt.Println("")

				found = true
				break
			}
		}

		if !found {
			req := &action.Pass{
				PlayerId: playerId,
			}
			cl.SendGameAction(action.TypePass, req)

			fmt.Printf("Pass\n")
			fmt.Println("")
		}

		// recv GameEvent
		for _, c := range clients {
			ge, _ := c.RecvGameEvent()
			if gameFinished(ge) {
				goto finish
			}
		}

		// abort
		if turn >= 100 {
			goto abort
		}
		turn++
	}
abort:
	assert.Failf(t, "aborted", "aborted")

finish:
}
