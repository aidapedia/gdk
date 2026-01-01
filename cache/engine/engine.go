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
	GET(ctx context.Context, key string) StringResult
	SET(ctx context.Context, key string, val interface{}, exp time.Duration) error
	HSET(ctx context.Context, key string, fields map[string]string) error
	HGET(ctx context.Context, key string, field string) StringResult
	HGETALL(ctx context.Context, key string) (map[string]string, error)
	DEL(ctx context.Context, keys ...string) error
}

// StringResult is the result of a cache operation.
type StringResult struct {
	value     string
	err       error
	unmarshal func(val string, dest interface{}) error
}

func (r StringResult) Value() string {
	return r.value
}

func (r StringResult) Err() error {
	return r.err
}

func (r StringResult) Scan(dest interface{}) error {
	return r.unmarshal(r.value, dest)
}
