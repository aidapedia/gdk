package secret

import (
	"context"
)

type Interface interface {
	GetSecret(ctx context.Context, target interface{}) error
}
