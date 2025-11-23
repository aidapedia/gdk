package engine

import (
	"context"
	"time"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

type GoRedisClient struct {
	*redis.Client
}

type GoRedisClientOpt struct {
	Opt               *redis.Options
	TracingInstrument bool
	MetricsInstrument bool
}

func NewGoRedisClient(opt GoRedisClientOpt) (Interface, error) {
	client := redis.NewClient(opt.Opt)
	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	if opt.TracingInstrument {
		if err := redisotel.InstrumentTracing(client); err != nil {
			return nil, err
		}
	}
	if opt.MetricsInstrument {
		if err := redisotel.InstrumentMetrics(client); err != nil {
			return nil, err
		}
	}

	return &GoRedisClient{
		Client: client,
	}, nil
}

func (c *GoRedisClient) GET(ctx context.Context, key string) (string, error) {
	return c.Client.Get(ctx, key).Result()
}

func (c *GoRedisClient) SET(ctx context.Context, key string, val interface{}, exp time.Duration) error {
	return c.Client.Set(ctx, key, val, exp).Err()
}

func (c *GoRedisClient) HSET(ctx context.Context, key string, fields map[string]string) error {
	values := []interface{}{}
	for field, val := range fields {
		values = append(values, field, val)
	}
	return c.Client.HSet(ctx, key, values...).Err()
}

func (c *GoRedisClient) HGET(ctx context.Context, key string, field string) (string, error) {
	return c.Client.HGet(ctx, key, field).Result()
}

func (c *GoRedisClient) HGETALL(ctx context.Context, key string) (map[string]string, error) {
	return c.Client.HGetAll(ctx, key).Result()
}

func (c *GoRedisClient) DEL(ctx context.Context, keys ...string) error {
	return c.Client.Del(ctx, keys...).Err()
}
