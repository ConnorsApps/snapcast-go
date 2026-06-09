package snapclient

import (
	"context"
	"fmt"
	"net/url"

	"github.com/coder/websocket"
)

func (c *Client) wsConnect(ctx context.Context) error {
	scheme := "ws"
	if c.secureConnection {
		scheme = "wss"
	}

	var (
		u   = url.URL{Scheme: scheme, Host: c.host, Path: "/jsonrpc"}
		err error
	)

	c.ws, _, err = websocket.Dial(ctx, u.String(), &websocket.DialOptions{
		HTTPClient: c.httpClient,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to snapcast at '%s', err: %w", c.host, err)
	}

	return nil
}

func (c *Client) Close() error {
	if err := c.ws.Close(websocket.StatusNormalClosure, ""); err != nil {
		return fmt.Errorf("failed to close snapcast websocket, err: %w", err)
	}
	return nil
}
