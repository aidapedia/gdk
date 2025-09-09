package callwrapper

import (
	"context"
	"fmt"
	"time"

	"github.com/aidapedia/gdk/callwrapper/pkg/cache"
	gdkLog "github.com/aidapedia/gdk/log"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
)

// CallWrapperName is the name of the callwrapper.
// This name will be used to generate the cache key.
// So make sure the name is unique.
var callWrappers = make(map[string]Interface)

// CallFunc is the function that will be called.
type callFunc = func(ctx context.Context) (interface{}, error)

// Interface is the interface for the callwrapper.
type Interface interface {
	Call(ctx context.Context, key map[string]interface{}, fn callFunc) (resp interface{}, err error)
}

// CallWrapper is the wrapper for the call function.
type CallWrapper struct {
	name  string
	sl    singleflight.Group
	cache cache.Client
	opt   Options
}

// Options is the configuration for the callwrapper
type Options struct {
	// Toggle for singleflight. Default is false.
	// Singleflight help to prevent multiple call to the same function.
	Singleflight bool

	// CacheExpiration is the expiration time for the cache. Default is 5 minutes.
	CacheExpiration time.Duration

	// Cache wll improve perfomance. But if you need realtime data response, set ths value to false
	Cache bool

	// CacheClient is the client for the cache.
	CacheClient cache.Client

	// Hook is the configuration for the hook.
	Hook Hook
}

// New creates a new callwrapper.
func New(name string, opt Options) error {
	cw := &CallWrapper{
		name: name,
		opt:  opt,
	}

	// Check if the callwrapper name already exists.
	if _, ok := callWrappers[name]; ok {
		return fmt.Errorf("callwrapper name %s already exists", name)
	}

	// Hook default configuration.
	if opt.Hook.OnErrorLog == nil {
		opt.Hook.OnErrorLog = func(ctx context.Context, msg string, err error) {
			gdkLog.ErrorCtx(ctx, msg, zap.Error(err))
		}
	}

	if opt.Hook.OnWarnLog == nil {
		opt.Hook.OnWarnLog = func(ctx context.Context, msg string, err error) {
			gdkLog.WarnCtx(ctx, msg, zap.Error(err))
		}
	}

	// Cache default configuration.
	if opt.Cache {
		if opt.CacheClient == nil {
			return fmt.Errorf("cache client is nil")
		}

		if opt.CacheExpiration == 0 {
			opt.CacheExpiration = time.Minute * 5
		}
	}

	callWrappers[name] = cw
	return nil
}

// GetCallWrapper returns the callwrapper by name.
// Be careful when using this function.
// If the callwrapper is not found, it will return nil.
func GetCallWrapper(name string) Interface {
	cw, ok := callWrappers[name]
	if !ok {
		return &CallWrapper{
			name: name,
		}
	}
	return cw
}

// Call executes the call function.
func (cw *CallWrapper) Call(ctx context.Context, key map[string]interface{}, fn callFunc) (resp interface{}, err error) {
	keyStr := generateKey(cw.name, key)
	if cw.opt.Cache {
		resp, err = cw.cache.Get(ctx, keyStr)
		if err != redis.Nil {
			return resp, nil
		}
	}

	defer func() {
		if err != nil {
			if cw.opt.Hook.OnFailure != nil {
				errHook := cw.opt.Hook.OnFailure(ctx)
				if errHook != nil {
					cw.opt.Hook.OnErrorLog(ctx, "failed to failure hook", errHook)
				}
			}
		} else {
			if cw.opt.Hook.OnSuccess != nil {
				errHook := cw.opt.Hook.OnSuccess(ctx)
				if errHook != nil {
					cw.opt.Hook.OnErrorLog(ctx, "failed to success hook", errHook)
				}
			}
		}
		if cw.opt.Cache {
			errCache := cw.cache.Set(ctx, keyStr, resp, cw.opt.CacheExpiration)
			if errCache != nil {
				cw.opt.Hook.OnWarnLog(ctx, "failed to set cache", err)
			}
		}
	}()

	if cw.opt.Singleflight {
		resp, err, _ = cw.sl.Do(keyStr, func() (interface{}, error) {
			return fn(ctx)
		})
		return resp, err
	}
	return fn(ctx)
}

func generateKey(name string, key map[string]interface{}) string {
	for k, v := range key {
		name += fmt.Sprintf(":%s:%s", k, v)
	}
	return name
}
