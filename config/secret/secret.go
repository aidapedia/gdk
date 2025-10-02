package secret

import (
	"context"
)

type Vault interface {
	GetSecret(ctx context.Context, target interface{}) error
}
