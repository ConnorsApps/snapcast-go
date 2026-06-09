package snapcast

import (
	"encoding/json"
	"time"
)

type RequestMethod string

// Requests
const (
	// Client
	MethodClientGetStatus  RequestMethod = "Client.GetStatus"
	MethodClientSetVolume  RequestMethod = "Client.SetVolume"
	MethodClientSetLatency RequestMethod = "Client.SetLatency"
	MethodClientSetName    RequestMethod = "Client.SetName"

	// Group
	MethodGroupGetStatus  RequestMethod = "Group.GetStatus"
	MethodGroupSetMute    RequestMethod = "Group.SetMute"
	MethodGroupSetStream  RequestMethod = "Group.SetStream"
	MethodGroupSetClients RequestMethod = "Group.SetClients"
	MethodGroupSetName    RequestMethod = "Group.SetName"

	// Server
	MethodServerGetRPCVersion RequestMethod = "Server.GetRPCVersion"
	MethodServerGetStatus     RequestMethod = "Server.GetStatus"
	MethodServerDeleteClient  RequestMethod = "Server.DeleteClient"

	// Stream
	MethodStreamAddStream    RequestMethod = "Stream.AddStream"
	MethodStreamRemoveStream RequestMethod = "Stream.RemoveStream"
	MethodStreamControl      RequestMethod = "Stream.Control"
	MethodStreamSetProperty  RequestMethod = "Stream.SetProperty"
)

type NotificationMethod string

// Notifications
const (
	// Client
	MethodClientOnConnect        NotificationMethod = "Client.OnConnect"
	MethodClientOnDisconnect     NotificationMethod = "Client.OnDisconnect"
	MethodClientOnVolumeChanged  NotificationMethod = "Client.OnVolumeChanged"
	MethodClientOnLatencyChanged NotificationMethod = "Client.OnLatencyChanged"
	MethodClientOnNameChanged    NotificationMethod = "Client.OnNameChanged"

	// Group
	MethodGroupOnMute          NotificationMethod = "Group.OnMute"
	MethodGroupOnStreamChanged NotificationMethod = "Group.OnStreamChanged"
	MethodGroupOnNameChanged   NotificationMethod = "Group.OnNameChanged"

	// Stream
	MethodStreamOnProperties NotificationMethod = "Stream.OnProperties"
	MethodStreamOnUpdate     NotificationMethod = "Stream.OnUpdate"

	// Server
	MethodServerOnUpdate NotificationMethod = "Server.OnUpdate"
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

type PlaybackStatus string

const (
	PlaybackPlaying PlaybackStatus = "playing"
	PlaybackPaused  PlaybackStatus = "paused"
	PlaybackStopped PlaybackStatus = "stopped"
)

type LoopStatus string

const (
	LoopNone     LoopStatus = "none"
	LoopTrack    LoopStatus = "track"
	LoopPlaylist LoopStatus = "playlist"
)

type StreamCommand string

const (
	StreamCommandNext        StreamCommand = "next"
	StreamCommandPrevious    StreamCommand = "previous"
	StreamCommandPlay        StreamCommand = "play"
	StreamCommandPause       StreamCommand = "pause"
	StreamCommandPlayPause   StreamCommand = "playPause"
	StreamCommandStop        StreamCommand = "stop"
	StreamCommandSeek        StreamCommand = "seek"
	StreamCommandSetPosition StreamCommand = "setPosition"
)

type (
	Error struct {
		Code    int         `json:"code,omitempty"`
		Data    interface{} `json:"data,omitempty"`
		Message string      `json:"message,omitempty"`
	}

	Notification struct {
		ID         *int                `json:"id,omitempty"`
		JsonRPC    string              `json:"jsonrpc,omitempty"`
		Method     *NotificationMethod `json:"method,omitempty"`
		Params     interface{}         `json:"params,omitempty"`
		ReceivedAt time.Time           `json:"-"`
	}

	Response struct {
		ID         *int        `json:"id,omitempty"`
		JsonRPC    string      `json:"jsonrpc,omitempty"`
		Error      *Error      `json:"error,omitempty"`
		Result     interface{} `json:"result,omitempty"`
		ReceivedAt time.Time   `json:"-"`
	}

	Request struct {
		ID      *int           `json:"id,omitempty"`
		JsonRPC string         `json:"jsonrpc,omitempty"`
		Method  *RequestMethod `json:"method,omitempty"`
		Params  interface{}    `json:"params,omitempty"`
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
		Properties *Properties `json:"properties,omitempty"`
	}

	Properties struct {
		PlaybackStatus PlaybackStatus `json:"playbackStatus,omitempty"`
		LoopStatus     LoopStatus     `json:"loopStatus,omitempty"`
		Shuffle        bool           `json:"shuffle,omitempty"`
		Volume         int            `json:"volume,omitempty"`
		Mute           bool           `json:"mute,omitempty"`
		Rate           float64        `json:"rate,omitempty"`
		Position       float64        `json:"position,omitempty"`
		CanGoNext      bool           `json:"canGoNext,omitempty"`
		CanGoPrevious  bool           `json:"canGoPrevious,omitempty"`
		CanPlay        bool           `json:"canPlay,omitempty"`
		CanPause       bool           `json:"canPause,omitempty"`
		CanSeek        bool           `json:"canSeek,omitempty"`
		CanControl     bool           `json:"canControl,omitempty"`
		Metadata       *Metadata      `json:"metadata,omitempty"`
	}

	Metadata struct {
		TrackID     string   `json:"trackId,omitempty"`
		Title       string   `json:"title,omitempty"`
		Artist      []string `json:"artist,omitempty"`
		Album       string   `json:"album,omitempty"`
		AlbumArtist []string `json:"albumArtist,omitempty"`
		Genre       []string `json:"genre,omitempty"`
		Duration    float64  `json:"duration,omitempty"`
		TrackNumber int      `json:"trackNumber,omitempty"`
		DiscNumber  int      `json:"discNumber,omitempty"`
		Date        string   `json:"date,omitempty"`
		ArtURL      string   `json:"artUrl,omitempty"`
		URL         string   `json:"url,omitempty"`
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
