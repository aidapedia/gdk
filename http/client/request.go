package client

import (
	"context"

	"github.com/gofiber/fiber/v3/client"
)

type Request struct {
	*client.Request
}

func NewRequest(ctx context.Context) *Request {
	req := client.AcquireRequest().SetContext(ctx)
	return &Request{
		Request: req,
	}
}
