package cache

import (
	"context"
	"time"
)

type Client interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, val interface{}, exp time.Duration) error
}
