package snapclient

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

func (n *Notifications) handleNotification(msg *snapcast.Notification) {
	switch *msg.Method {
	// --- Client
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

	case snapcast.MethodClientOnLatencyChanged:
		if n.ClientOnLatencyChanged == nil {
			return
		}
		var p = &snapcast.ClientOnLatencyChanged{}
		if err := marshalJSON(msg.Params, p); err != nil {
			n.readErr(err)
			return
		}
		n.ClientOnLatencyChanged <- p

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
	// --- Group
	case snapcast.MethodGroupOnMute:
		if n.GroupOnMute == nil {
			return
		}

		var p = &snapcast.GroupOnMute{}
		if err := marshalJSON(msg.Params, p); err != nil {
			n.readErr(err)
			return
		}
		n.GroupOnMute <- p

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

	case snapcast.MethodGroupOnNameChanged:
		if n.GroupOnNameChanged == nil {
			return
		}

		var p = &snapcast.GroupOnNameChanged{}
		if err := marshalJSON(msg.Params, p); err != nil {
			n.readErr(err)
			return
		}
		n.GroupOnNameChanged <- p
	// --- Stream
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

	case snapcast.MethodStreamOnProperties:
		if n.StreamOnProperties == nil {
			return
		}

		var p = &snapcast.StreamOnProperties{}
		if err := marshalJSON(msg.Params, p); err != nil {
			n.readErr(err)
			return
		}
		n.StreamOnProperties <- p

		// --- Server
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
	}
}
