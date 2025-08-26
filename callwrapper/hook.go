package callwrapper

import "context"

type Hook struct {
	// Hook that is called when success execution
	OnSuccess func(ctx context.Context) error

	// Hook that is called when falure execution
	OnFailure func(ctx context.Context) error

	// OnErrorLog is the hook that is called when error occurs
	OnErrorLog func(ctx context.Context, msg string, err error)

	// OnWarnLog is the hook that is called when error occurs but still tolerable.
	OnWarnLog func(ctx context.Context, msg string, err error)
}
