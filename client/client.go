package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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

type ClientOptions struct {
	Host             string
	RateLimiter      *rate.Limiter
	SecureConnection bool
}

func NewClient(o *ClientOptions) *Client {
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

func marshalJSON(i interface{}, to interface{}) error {
	raw, err := json.Marshal(i)
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, to)
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

func (c *Client) Listen(n *Notifications) {
	defer c.ws.Close()

	ch := make(chan *snapcast.Message, 5)

	go func() {
		for {
			_, raw, err := c.ws.ReadMessage()
			if err != nil {
				if strings.Contains(err.Error(), "close sent") {
					close(ch)
					// todo close out listener
					return
				} else {
					if n.MsgReaderErr != nil {
						n.MsgReaderErr <- err
					}
					continue
				}

			}
			var msg *snapcast.Message

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

	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				return
			}

			go n.handleMessage(msg)
		}
	}
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

		var p *snapcast.ServerOnUpdate
		if err := marshalJSON(msg.Params, p); err != nil {
			n.readErr(err)
			return
		}
		n.ServerOnUpdate <- p
	case snapcast.MethodStreamOnUpdate:
		if n.StreamOnUpdate == nil {
			return
		}

		var p *snapcast.StreamOnUpdate
		if err := marshalJSON(msg.Params, p); err != nil {
			n.readErr(err)
			return
		}
		n.StreamOnUpdate <- p

	case snapcast.MethodGroupOnStreamChanged:
		if n.GroupOnStreamChanged == nil {
			return
		}

		var p *snapcast.GroupOnStreamChanged
		if err := marshalJSON(msg.Params, p); err != nil {
			n.readErr(err)
			return
		}
		n.GroupOnStreamChanged <- p

	case snapcast.MethodClientOnConnect:
		if n.ClientOnConnect == nil {
			return
		}
		var p *snapcast.ClientOnConnect
		if err := marshalJSON(msg.Params, p); err != nil {
			n.readErr(err)
			return
		}
		n.ClientOnConnect <- p

	case snapcast.MethodClientOnDisconnect:
		if n.ClientOnDisconnect == nil {
			return
		}
		var p *snapcast.ClientOnDisconnect
		if err := marshalJSON(msg.Params, p); err != nil {
			n.readErr(err)
			return
		}
		n.ClientOnDisconnect <- p

	case snapcast.MethodClientOnVolumeChanged:
		if n.ClientOnVolumeChanged == nil {
			return
		}
		var p *snapcast.ClientOnVolumeChanged
		if err := marshalJSON(msg.Params, p); err != nil {
			n.readErr(err)
			return
		}
		n.ClientOnVolumeChanged <- p

	case snapcast.MethodClientOnNameChanged:
		if n.ClientOnNameChanged == nil {
			return
		}
		var p *snapcast.ClientOnNameChanged
		if err := marshalJSON(msg.Params, p); err != nil {
			n.readErr(err)
			return
		}
		n.ClientOnNameChanged <- p
	}
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
