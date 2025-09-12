package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aidapedia/gdk/callwrapper"
	"github.com/aidapedia/gdk/callwrapper/pkg/cache/goredis"
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
			}, func() (interface{}, error) {
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
			}, func() (interface{}, error) {
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
	}, func() (interface{}, error) {
		fmt.Println("calling GetUserByID")
		return repo.GetUserByID(ctx, 1)
	})
	if err != nil {
		panic(err)
	}
	user := resp.(User)
	fmt.Println("Without Singleflight User ID: ", user.ID)
}
