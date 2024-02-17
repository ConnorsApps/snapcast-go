package client

import (
	"fmt"
	"net/url"

	"github.com/gorilla/websocket"
)

func (c *Client) Connect() error {
	var (
		u   = url.URL{Scheme: "ws", Host: c.host, Path: "/jsonrpc"}
		err error
	)
	c.ws, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to dial to snapcast server at '%s', err: %w", c.host, err)
	}

	return nil
}
