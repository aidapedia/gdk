package client

import (
	"context"
	"fmt"
	"time"

	gctx "github.com/aidapedia/gdk/context"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"
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

// Send sends a request to the server.
// It sets the client and context to the request.
func (c *Client) Send(ctx context.Context, req *Request, resp interface{}) error {
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

	res, err := req.Send()
	if err != nil {
		return err
	}
	defer res.Close()

	if res.StatusCode() != 200 {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode())
	}

	err = c.cli.JSONUnmarshal()(res.Body(), resp)
	if err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}

func (c *Client) Get(ctx context.Context, req *Request, resp interface{}) error {
	req.SetMethod(fiber.MethodGet)
	return c.Send(ctx, req, resp)
}

func (c *Client) Post(ctx context.Context, req *Request, resp interface{}) error {
	req.SetMethod(fiber.MethodPost)
	return c.Send(ctx, req, resp)
}

func (c *Client) Put(ctx context.Context, req *Request, resp interface{}) error {
	req.SetMethod(fiber.MethodPut)
	return c.Send(ctx, req, resp)
}

func (c *Client) Delete(ctx context.Context, req *Request, resp interface{}) error {
	req.SetMethod(fiber.MethodDelete)
	return c.Send(ctx, req, resp)
}
