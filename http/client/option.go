package client

import (
	"time"

	"go.uber.org/ratelimit"
)

type Option interface {
	Apply(cli *Client)
}

// WithRateLimit is the option that adds rate limit to the client.
//
// Example:
//
//	WithRateLimit(1000)
func WithRateLimit(url string, rate int, opt ...ratelimit.Option) Option {
	return &withRateLimit{url: url, rate: rate, opt: opt}
}

type withRateLimit struct {
	url  string
	rate int
	opt  []ratelimit.Option
}

func (o *withRateLimit) Apply(cli *Client) {
	if o.rate > 0 {
		if cli.ratelimiter == nil {
			cli.ratelimiter = make(map[string]ratelimit.Limiter)
		}
		cli.ratelimiter[o.url] = ratelimit.New(o.rate, o.opt...)
	}
}

// WithGlobalTimeout is the option that adds global timeout to the client.
//
// Example:
//
//	WithGlobalTimeout(time.Second * 5)
func WithGlobalTimeout(timeout time.Duration) Option {
	return &withGlobalTimeout{timeout: timeout}
}

type withGlobalTimeout struct {
	timeout time.Duration
}

func (o *withGlobalTimeout) Apply(cli *Client) {
	cli.cli.SetTimeout(o.timeout)
}

func WithSetHeaders(headers map[string]string) Option {
	return &withSetHeaders{headers: headers}
}

type withSetHeaders struct {
	headers map[string]string
}

func (o *withSetHeaders) Apply(cli *Client) {
	cli.cli.SetHeaders(o.headers)
}
