package goredis

import (
	"context"
	"time"

	"github.com/aidapedia/gdk/callwrapper/pkg/cache"
	"github.com/go-redis/redis/v8"
)

type Client struct {
	*redis.Client
}

func New(opt *redis.Options) cache.Client {
	return &Client{
		Client: redis.NewClient(opt),
	}
}

func (c *Client) Get(ctx context.Context, key string) (interface{}, error) {
	return c.Client.Get(ctx, key).Result()
}
func (c *Client) Set(ctx context.Context, key string, val interface{}, exp time.Duration) error {
	return c.Client.Set(ctx, key, val, exp).Err()
}
