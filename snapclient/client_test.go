package snapclient

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ConnorsApps/snapcast-go/snapcast"
)

func TestClient(t *testing.T) {
	c := New(&Options{Host: "audio.connorskees.com:1780", SecureConnection: false})
	if _, err := c.Send(context.Background(), snapcast.MethodServerGetStatus, &snapcast.ServerGetStatusRequest{}); err != nil {
		t.Fatal(err)
	}

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

	if _, err := c.Listen(n); err != nil {
		t.Fatal(err)
	}

	defer c.Close()

	fmt.Println("Listening")

	time.Sleep(500 * time.Second)
}
