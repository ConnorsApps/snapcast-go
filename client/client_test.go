package client

import (
	"context"
	"testing"

	"github.com/ConnorsApps/snapcast-go/snapcast"
	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	a := assert.New(t)
	c := NewClient(&ClientOptions{Host: "audio.connorskees.com:1780", SecureConnection: false})
	_, err := c.Send(context.Background(), snapcast.ServerGetStatusMethod, &snapcast.ServerGetStatusRequest{})
	a.Nil(err)
}
