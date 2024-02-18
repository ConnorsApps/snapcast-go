package client

import (
	"encoding/json"

	"github.com/ConnorsApps/snapcast-go/snapcast"
)

func marshalJSON(i interface{}, to interface{}) error {
	raw, err := json.Marshal(i)
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, to)
}
func (n *Notifications) readErr(err error) {
	if n.MsgReaderErr != nil {
		n.MsgReaderErr <- err
	}
}
func (n *Notifications) handleMessage(msg *snapcast.Message) {
	switch *msg.Method {
	case snapcast.MethodServerOnUpdate:
		if n.ServerOnUpdate == nil {
			return
		}

		var p = &snapcast.ServerOnUpdate{}
		if err := marshalJSON(msg.Params, p); err != nil {
			n.readErr(err)
			return
		}
		n.ServerOnUpdate <- p
	case snapcast.MethodStreamOnUpdate:
		if n.StreamOnUpdate == nil {
			return
		}

		var p = &snapcast.StreamOnUpdate{}
		if err := marshalJSON(msg.Params, p); err != nil {
			n.readErr(err)
			return
		}
		n.StreamOnUpdate <- p

	case snapcast.MethodGroupOnStreamChanged:
		if n.GroupOnStreamChanged == nil {
			return
		}

		var p = &snapcast.GroupOnStreamChanged{}
		if err := marshalJSON(msg.Params, p); err != nil {
			n.readErr(err)
			return
		}
		n.GroupOnStreamChanged <- p

	case snapcast.MethodClientOnConnect:
		if n.ClientOnConnect == nil {
			return
		}
		var p = &snapcast.ClientOnConnect{}
		if err := marshalJSON(msg.Params, p); err != nil {
			n.readErr(err)
			return
		}
		n.ClientOnConnect <- p

	case snapcast.MethodClientOnDisconnect:
		if n.ClientOnDisconnect == nil {
			return
		}
		var p = &snapcast.ClientOnDisconnect{}
		if err := marshalJSON(msg.Params, p); err != nil {
			n.readErr(err)
			return
		}
		n.ClientOnDisconnect <- p

	case snapcast.MethodClientOnVolumeChanged:
		if n.ClientOnVolumeChanged == nil {
			return
		}
		var p = &snapcast.ClientOnVolumeChanged{}
		if err := marshalJSON(msg.Params, p); err != nil {
			n.readErr(err)
			return
		}
		n.ClientOnVolumeChanged <- p

	case snapcast.MethodClientOnNameChanged:
		if n.ClientOnNameChanged == nil {
			return
		}
		var p = &snapcast.ClientOnNameChanged{}
		if err := marshalJSON(msg.Params, p); err != nil {
			n.readErr(err)
			return
		}
		n.ClientOnNameChanged <- p
	}
}
