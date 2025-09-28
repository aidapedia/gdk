package engine

import (
	"context"
	"time"
)

type Field struct {
	Field string
	Value interface{}
}

type Interface interface {
	GET(ctx context.Context, key string) (string, error)
	SET(ctx context.Context, key string, val interface{}, exp time.Duration) error
	HSET(ctx context.Context, key string, fields map[string]string) error
	HGET(ctx context.Context, key string, field string) (string, error)
	HGETALL(ctx context.Context, key string) (map[string]string, error)
}
