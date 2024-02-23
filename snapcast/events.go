package snapcast

// Request Response
type (
	ClientSetVolumeRequest struct {
		ID     string `json:"id"`
		Volume Volume `json:"volume"`
	}

	ClientSetVolumeResponse struct {
		Volume Volume `json:"volume"`
	}

	GroupSetStreamRequest struct {
		ID       string `json:"id"`
		StreamID string `json:"stream_id"`
	}

	GroupSetStreamResponse struct {
		StreamID string `json:"stream_id"`
	}

	ServerGetStatusRequest struct{}

	ServerGetStatusResponse struct {
		Server Server `json:"server"`
	}
)

// Notifications
type (
	StreamOnUpdate struct {
		ID     string `json:"id"`
		Stream Stream `json:"stream"`
	}

	ServerOnUpdate struct {
		Server Server `json:"server"`
	}

	ClientOnConnect struct {
		Client *Client `json:"client"`
		ID     string  `json:"id"`
	}

	ClientOnDisconnect struct {
		Client *Client `json:"client"`
		ID     string  `json:"id"`
	}

	ClientOnLatencyChanged struct {
		Latency int    `json:"latency"`
		ID      string `json:"id"`
	}

	ClientOnVolumeChanged struct {
		Volume Volume `json:"volume"`
		ID     string `json:"id"`
	}

	ClientOnNameChanged struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	}

	GroupOnStreamChanged struct {
		ID       string `json:"id"`
		StreamId string `json:"stream_id"`
	}

	GroupOnNameChanged struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	}

	StreamOnProperties struct {
		ID       string            `json:"id"`
		Metadata map[string]string `json:"metadata"`
	}
)
