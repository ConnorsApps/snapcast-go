package snapcast

type ClientSetVolumeRequest struct {
	ID     string `json:"id"`
	Volume Volume `json:"volume"`
}

type ClientSetVolumeResponse struct {
	Volume Volume `json:"volume"`
}

type GroupSetStreamRequest struct {
	ID       string `json:"id"`
	StreamID string `json:"stream_id"`
}

type GroupSetStreamResponse struct {
	StreamID string `json:"stream_id"`
}

type ServerGetStatusRequest struct{}

type ServerGetStatusResponse struct {
	Server Server `json:"server"`
}

type StreamOnUpdate struct {
	ID     string `json:"id"`
	Stream Stream `json:"stream"`
}

type ServerOnUpdate struct {
	Server Server `json:"server"`
}

type ClientOnConnect struct {
	Client *Client `json:"client"`
	ID     string  `json:"id"`
}

type ClientOnDisconnect struct {
	Client *Client `json:"client"`
	ID     string  `json:"id"`
}

type ClientOnLatencyChanged struct {
	Latency int    `json:"latency"`
	ID      string `json:"id"`
}

type ClientOnVolumeChanged struct {
	Volume Volume `json:"volume"`
	ID     string `json:"id"`
}

type ClientOnNameChanged struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type GroupOnStreamChanged struct {
	ID       string `json:"id"`
	StreamId string `json:"stream_id"`
}

type GroupOnNameChanged struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type StreamOnProperties struct {
	ID       string            `json:"id"`
	Metadata map[string]string `json:"metadata"`
}
