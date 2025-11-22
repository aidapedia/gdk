package client

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3/client"
)

type Response struct {
	*client.Response
	err error
}

func NewResponse(res *client.Response, err error) *Response {
	return &Response{Response: res, err: err}
}

func (r *Response) Scan(v interface{}) error {
	return sonic.Unmarshal(r.Body(), v)
}

func (r *Response) Error() error {
	return r.err
}
