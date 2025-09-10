package callwrapper

import "context"

type Hook struct {
	// Hook that is called when success execution
	BeforeHook func(ctx context.Context) map[string]interface{}

	// Hook that is called when success execution
	AfterHook func(ctx context.Context, param map[string]interface{})

	// OnErrorLog is the hook that is called when error occurs
	OnErrorLog func(ctx context.Context, msg string, err error)

	// OnWarnLog is the hook that is called when error occurs but still tolerable.
	OnWarnLog func(ctx context.Context, msg string, err error)
}
