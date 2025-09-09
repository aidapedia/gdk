package gocache

import (
	"context"
	"errors"
	"time"

	"github.com/aidapedia/gdk/callwrapper/pkg/cache"
	memcache "github.com/patrickmn/go-cache"
)

type Client struct {
	cache *memcache.Cache
	redis cache.Client
}

func New(exp, purge time.Duration, redis cache.Client) cache.Client {
	return &Client{
		cache: memcache.New(exp, purge),
		redis: redis,
	}
}

func (c *Client) Get(ctx context.Context, key string) (interface{}, error) {
	resp, isSuccess := c.cache.Get(key)
	if !isSuccess {
		if c.redis == nil {
			return nil, errors.New("cache missed")
		}
		return c.redis.Get(ctx, key)
	}
	return resp, nil
}

func (c *Client) Set(ctx context.Context, key string, val interface{}, exp time.Duration) error {
	c.cache.Set(key, val, exp)
	return c.redis.Set(ctx, key, val, exp)
}
