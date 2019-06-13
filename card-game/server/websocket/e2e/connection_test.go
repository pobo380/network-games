package e2e

import (
	"testing"
)

func Test_Connect(t *testing.T) {
	c := newWssConnection()
	defer c.Con.Close()
}

func Test_Reconnect(t *testing.T) {
	c := newWssConnection()
	defer c.Con.Close()

	c = newWssConnectionWithArgs(DefaultWssEndpoint, c.PlayerId)
	defer c.Con.Close()
}
