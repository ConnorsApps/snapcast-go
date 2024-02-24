package snapclient

import (
	"encoding/json"

	"github.com/ConnorsApps/snapcast-go/snapcast"
	"github.com/gorilla/websocket"
)

func (c *Client) readMessages(n *Notifications, msgChan chan *snapcast.Message, wsClose chan error) {
	for {
		_, raw, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				close(msgChan)
				wsClose <- err
				return
			}

			if n.MsgReaderErr != nil {
				n.MsgReaderErr <- err
			}
			continue
		}

		var msg = &snapcast.Message{}

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
