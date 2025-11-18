package module

import (
	"context"
	"fmt"
)

type Module string

const (
	FileModule   Module = "file"
	ConsulModule Module = "consul"
)

// ErrKeyNotFound is the error when key is not found.
var ErrKeyNotFound = fmt.Errorf("key not found")

type Interface interface {
	// GetValue returns the value of the key.
	GetValue(ctx context.Context, key string) (interface{}, error)

	GetBool(ctx context.Context, key string) (bool, error)
	GetInt(ctx context.Context, key string) (int, error)
	GetString(ctx context.Context, key string) (string, error)
	GetStruct(ctx context.Context, key string, v interface{}) error

	Watch() (chan bool, error)
}
