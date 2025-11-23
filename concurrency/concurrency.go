package concurrency

import (
	"context"
	"fmt"

	gerr "github.com/aidapedia/gdk/error"
	"github.com/aidapedia/gdk/log"
	"go.uber.org/zap"
)

var concurrency routine

type routine struct {
	recoverHook func(ctx context.Context, err interface{})
}

func init() {
	concurrency = routine{
		recoverHook: func(ctx context.Context, err interface{}) {
			log.ErrorCtx(ctx, "routine panic", zap.Any("error", err))
		},
	}
}

// SetRecoverHook sets the recover hook for the routine.
//
// Custom recover hook can be used to handle panic in routine. For example: send notification
func SetRecoverHook(fn func(ctx context.Context, err interface{})) {
	concurrency.recoverHook = fn
}

// Call runs the function in a new routine.
//
// It is useful when you want to run a function in a new routine.
func Call(ctx context.Context, fn func(ctx context.Context)) {
	ctxRtn := context.WithoutCancel(ctx)
	// Capture panic to recover
	go func() {
		defer func() {
			if r := recover(); r != nil {
				concurrency.recoverHook(ctxRtn, gerr.New(fmt.Errorf("%v", r)))
			}
		}()
		fn(ctxRtn)
	}()
}
