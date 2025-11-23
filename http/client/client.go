package client

import (
	"context"

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
}

// New creates a new client.
func New(opt ...Option) *Client {
	c := &Client{}
	// create client
	cli := client.New()
	cli.SetJSONUnmarshal(sonic.Unmarshal)
	c.cli = cli

	for _, o := range opt {
		o.Apply(c)
	}

	return c
}

func (c *Client) Send(ctx context.Context, req *Request) (*client.Response, error) {
	// force set client to the request
	req.SetClient(c.cli).AddHeader(gctx.ContextKeyLogID, ctx.Value(gctx.ContextKeyLogID).(string))

	// check rate limit
	if c.ratelimiter != nil {
		if limiter, ok := c.ratelimiter[req.URL()]; ok {
			limiter.Take()
		}
	}

	return req.Send()
}
