package client

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ConnorsApps/snapcast-go/snapcast"
	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	a := assert.New(t)
	c := NewClient(&ClientOptions{Host: "audio.connorskees.com:1780", SecureConnection: false})
	_, err := c.Send(context.Background(), snapcast.MethodServerGetStatus, &snapcast.ServerGetStatusRequest{})
	a.Nil(err)

	onUpdate := make(chan *snapcast.ClientOnVolumeChanged)
	n := &Notifications{
		ClientOnVolumeChanged: onUpdate,
	}

	go func() {
		for {
			msg := <-onUpdate
			fmt.Println("onUpdate", msg)
		}
	}()

	cancel, err := c.Listen(n)
	defer cancel()
	a.Nil(err)

	fmt.Println("Listening")

	time.Sleep(500 * time.Second)
}
