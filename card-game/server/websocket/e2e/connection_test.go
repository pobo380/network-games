package e2e

import (
	"testing"
)

func Test_Connect(t *testing.T) {
	con, _ := newWssConnection()
	defer con.Close()
}

func Test_Reconnect(t *testing.T) {
	con, playerId := newWssConnection()
	con.Close()

	con = newWssConnectionWithArgs(DefaultWssEndpoint, playerId)
	defer con.Close()
}
