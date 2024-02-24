package snapcast

// Request Response
type (
	ClientGetStatusRequest struct {
		ID string `json:"id"`
	}

	ClientGetStatusResponse struct {
		Client Client `json:"client"`
	}

	ClientSetVolumeRequest struct {
		ID     string `json:"id"`
		Volume Volume `json:"volume"`
	}

	ClientSetVolumeResponse struct {
		Volume Volume `json:"volume"`
	}

	ClientSetLatencyRequest struct {
		ID      string `json:"id"`
		Latency int    `json:"latency"`
	}

	ClientSetLatencyResponse struct {
		Latency int `json:"latency"`
	}

	ClientSetNameRequest struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	ClientSetNameResponse struct {
		Name string `json:"name"`
	}

	GroupGetStatusRequest struct {
		ID string `json:"id"`
	}

	GroupGetStatusResponse struct {
		Group Group `json:"group"`
	}

	GroupSetMuteRequest struct {
		ID    string `json:"id"`
		Muted bool   `json:"muted"`
	}

	GroupSetMuteResponse struct {
		Muted bool `json:"muted"`
	}

	GroupSetStreamRequest struct {
		ID       string `json:"id"`
		StreamID string `json:"stream_id"`
	}

	GroupSetStreamResponse struct {
		StreamID string `json:"stream_id"`
	}

	GroupSetClientsRequest struct {
		ID      string   `json:"id"`
		Clients []string `json:"clients"`
	}

	GroupSetClientsResponse struct {
		Clients []*Client `json:"clients"`
	}

	GroupSetNameRequest struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	GroupSetNameResponse struct {
		Name string `json:"name"`
	}

	ServerGetRPCVersion struct{}

	ServerGetRPCVersionResponse struct {
		Major int `json:"major"`
		Minor int `json:"minor"`
		Patch int `json:"patch"`
	}

	ServerGetStatusRequest struct{}

	ServerGetStatusResponse struct {
		Server Server `json:"server"`
	}

	ServerDeleteClient struct {
		ID string `json:"id"`
	}

	ServerDeleteClientResponse struct {
		Server Server `json:"server"`
	}

	StreamAddStream struct {
		StreamUri string `json:"streamUri"`
	}

	StreamAddStreamResponse struct {
		StreamId string `json:"stream_id"`
	}

	StreamRemoveStream struct {
		ID string `json:"id"`
	}

	StreamRemoveStreamResponse struct {
		StreamId string `json:"stream_id"`
	}

	StreamControl struct {
		Command string      `json:"command"`
		Params  interface{} `json:"params"`
	}

	StreamControlResponse string

	StreamSetPropety struct {
		ID       string      `json:"id"`
		Property interface{} `json:"property"`
		Value    interface{} `json:"value"`
	}

	StreamSetPropetyResponse string
)

// Notifications
type (
	ClientOnConnect struct {
		Client *Client `json:"client"`
		ID     string  `json:"id"`
	}

	ClientOnDisconnect struct {
		Client *Client `json:"client"`
		ID     string  `json:"id"`
	}

	ClientOnVolumeChanged struct {
		Volume Volume `json:"volume"`
		ID     string `json:"id"`
	}

	ClientOnLatencyChanged struct {
		Latency int    `json:"latency"`
		ID      string `json:"id"`
	}

	ClientOnNameChanged struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	}

	GroupOnMute struct {
		Mute bool   `json:"mute"`
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

	StreamOnUpdate struct {
		ID     string `json:"id"`
		Stream Stream `json:"stream"`
	}

	StreamOnProperties struct {
		ID       string            `json:"id"`
		Metadata map[string]string `json:"metadata"`
	}

	ServerOnUpdate struct {
		Server Server `json:"server"`
	}
)
