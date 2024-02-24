package snapclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/ConnorsApps/snapcast-go/snapcast"
	"github.com/gorilla/websocket"
	"golang.org/x/time/rate"
)

var (
	DefaultRequestBurst = 10
	DefaultRequestRate  = rate.Every(time.Second / 2)
)

type state struct {
	sync.Mutex
	reqCount uint
}

type Client struct {
	limiter          *rate.Limiter
	ws               *websocket.Conn
	host             string
	state            state
	secureConnection bool
	httpClient       *http.Client
}

type Options struct {
	Host        string
	RateLimiter *rate.Limiter
	// if secure then https & wss and used else http & ws protocols
	SecureConnection bool
}

func New(o *Options) *Client {
	if o.RateLimiter == nil {
		o.RateLimiter = rate.NewLimiter(3, DefaultRequestBurst)
	}

	return &Client{
		host:    o.Host,
		limiter: o.RateLimiter,
		httpClient: &http.Client{Transport: &http.Transport{
			MaxIdleConnsPerHost: -1, // Todo, better handle keep alive between snapcast and client
		}},
		secureConnection: o.SecureConnection,
	}
}

type Notifications struct {
	MsgReaderErr chan error

	// Client
	ClientOnConnect        chan *snapcast.ClientOnConnect
	ClientOnDisconnect     chan *snapcast.ClientOnDisconnect
	ClientOnVolumeChanged  chan *snapcast.ClientOnVolumeChanged
	ClientOnLatencyChanged chan *snapcast.ClientOnLatencyChanged
	ClientOnNameChanged    chan *snapcast.ClientOnNameChanged
	// Group
	GroupOnMute          chan *snapcast.GroupOnMute
	GroupOnStreamChanged chan *snapcast.GroupOnStreamChanged
	GroupOnNameChanged   chan *snapcast.GroupOnNameChanged
	// Stream
	StreamOnUpdate     chan *snapcast.StreamOnUpdate
	StreamOnProperties chan *snapcast.StreamOnProperties
	// Server
	ServerOnUpdate chan *snapcast.ServerOnUpdate
}

// Passes a websocket closer channel or an error on initial setup
func (c *Client) Listen(n *Notifications) (chan error, error) {
	var (
		msgChan = make(chan *snapcast.Notification, 5)
		wsClose = make(chan error)
	)

	wsClose = make(chan error)

	if err := c.wsConnect(); err != nil {
		return wsClose, err
	}
	go c.readNotifications(n, msgChan, wsClose)

	go func() {
		for {
			msg, ok := <-msgChan
			if !ok { // Websocket closed
				return
			}

			n.handleNotification(msg)
		}
	}()

	return wsClose, nil
}

func (c *Client) Send(ctx context.Context, method snapcast.RequestMethod, params interface{}) (*snapcast.Response, error) {
	c.state.Lock()
	c.state.reqCount += 1
	c.state.Unlock()

	var (
		id  = int(c.state.reqCount)
		req = snapcast.Request{
			ID:      &id,
			JsonRPC: "2.0",
			Method:  &method,
			Params:  &params,
		}
	)

	var response = &snapcast.Response{}

	// Limit requests/sec so we don't DOS the poor server
	if err := c.limiter.Wait(ctx); err != nil {
		return response, err
	}

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(req); err != nil {
		return response, fmt.Errorf("json.NewEncoder(buf).Encode(req): %v", err)
	}

	proto := "http"
	if c.secureConnection {
		proto = "https"
	}
	var url = proto + "://" + c.host + "/jsonrpc"

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, buf)
	if err != nil {
		return response, fmt.Errorf("http.NewRequestWithContext: %v", err)
	}

	httpReq.Header = http.Header{
		"Accept":       {"application/json"},
		"Content-Type": {"application/json"},
	}

	res, err := c.httpClient.Do(httpReq)
	if err != nil {
		return response, err
	}

	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return response, fmt.Errorf("%s", res.Status)
	}

	if err := json.NewDecoder(res.Body).Decode(response); err != nil {
		return response, fmt.Errorf("json.NewDecoder: %v", err)
	}

	return response, nil
}
