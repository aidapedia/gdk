package client

import (
	"context"
	"fmt"
	"time"

	gctx "github.com/aidapedia/gdk/context"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3/client"
	"go.uber.org/ratelimit"
)

// Client is the client for HTTP requests.
type Client struct {
	cli *client.Client
	// rate limiter
	// Key is the path of the request. Example: using path /user
	ratelimiter map[string]ratelimit.Limiter
	// global timeout for all requests
	globalTimeout time.Duration
}

// New creates a new client.
func New(opt ...Option) *Client {
	c := &Client{}
	for _, o := range opt {
		o.Apply(c)
	}
	// create client
	cli := client.New()
	cli.SetJSONUnmarshal(sonic.Unmarshal)
	c.cli = cli

	return c
}

func (c *Client) Send(ctx context.Context, req *Request) *Response {
	// force set client to the request
	req.SetClient(c.cli)
	req.AddHeader(gctx.ContextKeyLogID, ctx.Value(gctx.ContextKeyLogID).(string))

	// check rate limit
	if c.ratelimiter != nil {
		if limiter, ok := c.ratelimiter[req.URL()]; ok {
			limiter.Take()
		}
	}

	// set global timeout if not set
	// priority: request timeout > client global timeout
	if req.Timeout() == 0 && c.globalTimeout > 0 {
		req.SetTimeout(c.globalTimeout)
	}

	resp, err := req.Send()
	if err != nil {
		return NewResponse(nil, err)
	}
	defer resp.Close()

	// check status code
	if resp.StatusCode() >= 400 {
		return NewResponse(resp, fmt.Errorf("unexpected status code: %d", resp.StatusCode()))
	}

	return NewResponse(resp, nil)
}
