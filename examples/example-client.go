package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ConnorsApps/snapcast-go/snapcast"
	"github.com/ConnorsApps/snapcast-go/snapclient"
)

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	var client = snapclient.New(&snapclient.Options{
		Host:             "localhost:1780",
		SecureConnection: false,
	})

	var notify = &snapclient.Notifications{
		MsgReaderErr:          make(chan error),
		ServerOnUpdate:        make(chan *snapcast.ServerOnUpdate),
		ClientOnVolumeChanged: make(chan *snapcast.ClientOnVolumeChanged),
		ClientOnNameChanged:   make(chan *snapcast.ClientOnNameChanged),
	}

	wsClose, err := client.Listen(notify)
	check(err)

	// Listen for events
	go func() {
		for {
			select {
			case m := <-notify.MsgReaderErr:
				fmt.Println("Message reader error", m)
				continue
			case m := <-notify.ServerOnUpdate:
				fmt.Println("ServerOnUpdate", m)
			case m := <-notify.ClientOnVolumeChanged:
				fmt.Println("ClientOnVolumeChanged", m)
			case m := <-notify.ClientOnNameChanged:
				fmt.Println("ClientOnNameChanged", m)
			}
		}
	}()

	res, err := client.Send(context.Background(), snapcast.MethodServerGetStatus, struct{}{})
	check(err)
	if res.Error != nil {
		log.Fatalln(res.Error)
	}

	initialState, err := snapcast.ParseResult[snapcast.ServerGetStatusResponse](res.Result)
	check(err)
	fmt.Println("Initial state", initialState)

	panic(<-wsClose)
}
