package client

import (
	"context"

	"github.com/gofiber/fiber/v3/client"
)

type Client struct {
	cli *client.Client
}

func New() *Client {
	return &Client{
		cli: client.New(),
	}
}

// Send sends a request to the server.
// It sets the client and context to the request.
func (c *Client) Send(ctx context.Context, req *client.Request) (*client.Response, error) {
	req.SetClient(c.cli)
	req.SetContext(ctx)
	return req.Send()
}
