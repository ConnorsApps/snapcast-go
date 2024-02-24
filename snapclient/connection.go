package snapclient

import (
	"fmt"
	"net/url"

	"github.com/gorilla/websocket"
)

func (c *Client) wsConnect() error {
	scheme := "ws"
	if c.secureConnection {
		scheme = "wss"
	}

	var (
		u   = url.URL{Scheme: scheme, Host: c.host, Path: "/jsonrpc"}
		err error
	)

	c.ws, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to connect to snapcast at '%s', err: %w", c.host, err)
	}

	return nil
}

func (c *Client) Close() error {
	err := c.ws.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
	)
	if err != nil {
		return fmt.Errorf("failed to write close message to snapcast, err: %w", err)
	}

	return c.ws.Close()
}
