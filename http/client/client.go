package client

import (
	"context"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/client"
)

type Client struct {
	cli *client.Client
}

func New() *Client {
	cli := client.New()
	cli.SetJSONUnmarshal(sonic.Unmarshal)
	return &Client{
		cli: cli,
	}
}

// Send sends a request to the server.
// It sets the client and context to the request.
func (c *Client) send(ctx context.Context, req *Request, resp interface{}) error {
	if req.Client() == nil {
		req.SetClient(c.cli)
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
	return c.send(ctx, req, resp)
}

func (c *Client) Post(ctx context.Context, req *Request, resp interface{}) error {
	req.SetMethod(fiber.MethodPost)
	return c.send(ctx, req, resp)
}

func (c *Client) Put(ctx context.Context, req *Request, resp interface{}) error {
	req.SetMethod(fiber.MethodPut)
	return c.send(ctx, req, resp)
}

func (c *Client) Delete(ctx context.Context, req *Request, resp interface{}) error {
	req.SetMethod(fiber.MethodDelete)
	return c.send(ctx, req, resp)
}
