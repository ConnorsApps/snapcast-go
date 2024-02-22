package snapcast

import (
	"encoding/json"
	"time"
)

type Method string

const (
	MethodClientSetVolume Method = "Client.SetVolume"
	MethodGroupSetStream  Method = "Group.SetStream"
	MethodServerGetStatus Method = "Server.GetStatus"

	MethodServerOnUpdate        Method = "Server.OnUpdate"
	MethodStreamOnUpdate        Method = "Stream.OnUpdate"
	MethodGroupOnStreamChanged  Method = "Group.OnStreamChanged"
	MethodClientOnConnect       Method = "Client.OnConnect"
	MethodClientOnDisconnect    Method = "Client.OnDisconnect"
	MethodClientOnVolumeChanged Method = "Client.OnVolumeChanged"
	MethodClientOnNameChanged   Method = "Client.OnNameChanged"
)

type StreamStatus string

const (
	StreamIdle    StreamStatus = "idle"
	StreamPlaying StreamStatus = "playing"
)

func (s StreamStatus) IsPlaying() bool {
	return s == StreamPlaying
}

func (s StreamStatus) IsIdle() bool {
	return s == StreamIdle
}

type (
	Error struct {
		Code    int         `json:"code,omitempty"`
		Data    interface{} `json:"data,omitempty"`
		Message string      `json:"message,omitempty"`
	}

	Message struct {
		ID         *int        `json:"id,omitempty"`
		JsonRPC    string      `json:"jsonrpc,omitempty"`
		Method     *Method     `json:"method,omitempty"`
		Error      *Error      `json:"error,omitempty"`
		Result     interface{} `json:"result,omitempty"`
		Params     interface{} `json:"params,omitempty"`
		ReceivedAt time.Time   `json:"-"`
	}
)

func ParseResult[T any](result interface{}) (*T, error) {
	var t = new(T)
	data, err := json.Marshal(result)
	if err != nil {
		return t, err
	}

	err = json.Unmarshal(data, t)

	return t, err
}

type (
	Stream struct {
		ID     string       `json:"id"`
		Status StreamStatus `json:"status"`
		URI    struct {
			Fragment string            `json:"fragment"`
			Host     string            `json:"host"`
			Path     string            `json:"path"`
			Query    map[string]string `json:"query"`
			Raw      string            `json:"raw"`
			Scheme   string            `json:"scheme"`
		} `json:"uri"`
	}

	Snapserver struct {
		ControlProtocolVersion int    `json:"controlProtocolVersion"`
		Name                   string `json:"name"`
		ProtocolVersion        int    `json:"protocolVersion"`
		Version                string `json:"version"`
	}

	Volume struct {
		Muted   bool `json:"muted"`
		Percent int  `json:"percent"`
	}

	Client struct {
		Config struct {
			Instance int    `json:"instance"`
			Latency  int    `json:"latency"`
			Name     string `json:"name"`
			Volume   Volume `json:"volume"`
		} `json:"config"`
		Connected bool `json:"connected"`
		Host      struct {
			Arch string `json:"arch"`
			IP   string `json:"ip"`
			MAC  string `json:"mac"`
			Name string `json:"name"`
			OS   string `json:"os"`
		} `json:"host"`
		ID       string `json:"id"`
		LastSeen struct {
			Sec  int `json:"sec"`
			USec int `json:"usec"`
		} `json:"lastSeen"`
		Snapclient struct {
			Name            string `json:"name"`
			ProtocolVersion int    `json:"protocolVersion"`
			Version         string `json:"version"`
		} `json:"snapclient"`
	}

	Group struct {
		Clients  []Client `json:"clients"`
		ID       string   `json:"id"`
		Muted    bool     `json:"muted"`
		Name     string   `json:"name"`
		StreamID string   `json:"stream_id"`
	}

	Server struct {
		Groups []Group `json:"groups"`
		Host   struct {
			Arch string `json:"arch"`
			IP   string `json:"ip"`
			MAC  string `json:"mac"`
			Name string `json:"name"`
			OS   string `json:"os"`
		} `json:"host"`
		Snapserver Snapserver `json:"snapserver"`
		Streams    []Stream   `json:"streams"`
	}
)
