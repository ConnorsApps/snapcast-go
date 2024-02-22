package snapclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/ConnorsApps/snapcast-go/snapcast"
	"github.com/gorilla/websocket"
	"golang.org/x/time/rate"
)

var (
	DefaultRequestBurst = 10
	DefaultRate         = rate.Every(time.Second / 2)
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
		o.RateLimiter = rate.NewLimiter(DefaultRate, DefaultRequestBurst)
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
	MsgReaderErr          chan error
	ServerOnUpdate        chan *snapcast.ServerOnUpdate
	StreamOnUpdate        chan *snapcast.StreamOnUpdate
	GroupOnStreamChanged  chan *snapcast.GroupOnStreamChanged
	ClientOnConnect       chan *snapcast.ClientOnConnect
	ClientOnDisconnect    chan *snapcast.ClientOnDisconnect
	ClientOnVolumeChanged chan *snapcast.ClientOnVolumeChanged
	ClientOnNameChanged   chan *snapcast.ClientOnNameChanged
}

func (c *Client) wsConnect() error {
	scheme := "ws"
	if c.secureConnection {
		scheme = "wss"
	}

	var (
		u   = url.URL{Scheme: scheme, Host: c.host, Path: "/jsonrpc"}
		err error
	)

	c.ws, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to connect to snapcast at '%s', err: %w", c.host, err)
	}

	return nil
}

func (c *Client) Close() error {
	err := c.ws.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
	)
	if err != nil {
		return fmt.Errorf("failed to write close message to snapcast, err: %w", err)
	}

	c.ws.Close()
	return nil
}

// Passes a websocket closer channel or an error on initial setup
func (c *Client) Listen(n *Notifications) (chan error, error) {
	var (
		ch      = make(chan *snapcast.Message, 5)
		wsClose = make(chan error)
	)

	if err := c.wsConnect(); err != nil {
		return wsClose, err
	}

	go func() {
		for {
			_, raw, err := c.ws.ReadMessage()
			if err != nil {
				fmt.Printf("err %v \n", err)

				if strings.Contains(err.Error(), "(abnormal closure): unexpected EOF") {
					close(ch)
					wsClose <- nil
					return
				} else if strings.Contains(err.Error(), "close sent") {
					close(ch)
					wsClose <- nil
					return
				} else {
					if n.MsgReaderErr != nil {
						wsClose <- err
					}
					continue
				}

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
				ch <- msg
			}
		}
	}()

	go func() {
		for {
			msg, ok := <-ch
			if !ok { // Websocket closed
				return
			}

			go n.handleMessage(msg)
		}
	}()

	return wsClose, nil
}

func (c *Client) Send(ctx context.Context, method snapcast.Method, params interface{}) (*snapcast.Message, error) {
	c.state.Lock()
	c.state.reqCount += 1
	c.state.Unlock()

	var (
		id  = int(c.state.reqCount)
		req = snapcast.Message{
			ID:      &id,
			JsonRPC: "2.0",
			Method:  &method,
			Params:  &params,
		}
	)

	var response = &snapcast.Message{}

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
