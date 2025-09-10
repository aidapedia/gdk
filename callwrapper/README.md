# Callwrapper

Callwrapper is a package that helps you to wrap your function with a middleware.

## Features
- <b>Caching support</b>. Cache the result of your function. This will help to improve latency. Help to reduce the load on your external service.
- <b>Customize Caching Client</b>. You can customize the caching client. We have provide some example on ```pkg/cache```.
- <b>Singleflight support</b>. You can use singleflight to prevent multiple requests from hitting the same function.
- <b>Hook support</b>. You can add hook to your function. This will help you to add some extra logic to your function.
	- <b>Before Hook</b>. This hook will be called before the function is called.
	- <b>After Hook</b>. This hook will be called after the function is called.
	- <b>OnErrorLog</b>. This hook will be called when the function returns an error.
    - <b>OnWarnLog</b>. This hook will be called when the function returns a warning.

## Roadmap
- <b>Support Circuit Breaker</b>. You can use circuit breaker to prevent your service to call external service when the external service is down.
- <b>Support Telemetry</b>. Help to monitoring your external call like latency, error rate and QPS.

## How to use
This is an example of how to use callwrapper.
```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aidapedia/callwrapper"
	"github.com/aidapedia/callwrapper/pkg/cache/goredis"
	"github.com/go-redis/redis/v8"
)

type User struct {
	ID int
}

type Repository struct {
	// do external call
}

func NewRepository() *Repository {
	// add new callwrapper
	redisClient := goredis.New(&redis.Options{
		Addr: "localhost:6379",
	})

	callwrapper.New("GetUserByID", callwrapper.Options{})
	// with redis cache
	callwrapper.New("GetUserByIDwithRedis", callwrapper.Options{
		Cache:       true,
		CacheClient: redisClient,
	})
	// with singleflight
	callwrapper.New("GetUserByIDwithSL", callwrapper.Options{
		Singleflight: true,
	})

	return &Repository{}
}

func (r *Repository) GetUserByID(ctx context.Context, id int) (User, error) {
	time.Sleep(1 * time.Second)
	return User{
		ID: id,
	}, nil
}

func main() {
	repo := NewRepository()
	ctx := context.Background()

	// with singleflight
	for i := 0; i < 2; i++ {
		var resp interface{}
		var err error
		go func() {
			resp, err = callwrapper.GetCallWrapper("GetUserByIDwithSL").Call(ctx, map[string]interface{}{
				"id": 1,
			}, func(ctx context.Context) (interface{}, error) {
				fmt.Println("calling GetUserByIDwithSL")
				return repo.GetUserByID(ctx, 1)
			})
			if err != nil {
				panic(err)
			}
			user := resp.(User)
			fmt.Println("Singleflight User ID: ", user.ID)
		}()
	}
	time.Sleep(3 * time.Second)

	// with redis cache
	for i := 0; i < 2; i++ {
		var resp interface{}
		var err error
		go func() {
			resp, err = callwrapper.GetCallWrapper("GetUserByIDwithRedis").Call(ctx, map[string]interface{}{
				"id": 1,
			}, func(ctx context.Context) (interface{}, error) {
				fmt.Println("calling GetUserByIDwithRedis")
				return repo.GetUserByID(ctx, 1)
			})
			if err != nil {
				panic(err)
			}
			user := resp.(User)
			fmt.Println("Redis User ID: ", user.ID)
		}()
	}
	time.Sleep(3 * time.Second)

	// without singleflight
	resp, err := callwrapper.GetCallWrapper("GetUserByID").Call(ctx, map[string]interface{}{
		"id": 1,
	}, func(ctx context.Context) (interface{}, error) {
		fmt.Println("calling GetUserByID")
		return repo.GetUserByID(ctx, 1)
	})
	if err != nil {
		panic(err)
	}
	user := resp.(User)
	fmt.Println("Without Singleflight User ID: ", user.ID)
}

```
