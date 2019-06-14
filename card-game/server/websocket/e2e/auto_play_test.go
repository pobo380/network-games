package e2e

import (
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
	cA.RecvRoomInfo()
	cB.RecvRoomInfo()
	cC.RecvRoomInfo()
	cD.RecvRoomInfo()

	// recv GameStart
	cA.RecvGameStart()
	cB.RecvGameStart()
	cC.RecvGameStart()
	cD.RecvGameStart()

	// recv GameEvent
	cA.RecvGameEvent()
	cB.RecvGameEvent()
	cC.RecvGameEvent()
	cD.RecvGameEvent()

	assert.NotNil(t, cA.State)
	assert.NotNil(t, cB.State)
	assert.NotNil(t, cC.State)
	assert.NotNil(t, cD.State)
}
