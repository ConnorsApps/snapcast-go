package snapclient

import (
	"context"
	"encoding/json"

	"github.com/ConnorsApps/snapcast-go/snapcast"
	"github.com/coder/websocket"
)

func (c *Client) readNotifications(ctx context.Context, n *Notifications, msgChan chan *snapcast.Notification, wsClose chan error) {
	for {
		_, raw, err := c.ws.Read(ctx)
		if err != nil {
			if status := websocket.CloseStatus(err); status != -1 && status != websocket.StatusNormalClosure {
				close(msgChan)
				wsClose <- err
				return
			}

			if ctx.Err() != nil {
				close(msgChan)
				wsClose <- err
				return
			}

			if n.MsgReaderErr != nil {
				n.MsgReaderErr <- err
			}
			continue
		}

		var msg = &snapcast.Notification{}

		if err := json.Unmarshal(raw, msg); err != nil {
			if n.MsgReaderErr != nil {
				n.MsgReaderErr <- err
			}
			continue
		}

		// Only process notifications, not responses to requests
		if msg.Params != nil && msg.Method != nil {
			msgChan <- msg
		}
	}
}
