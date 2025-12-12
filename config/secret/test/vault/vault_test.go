package vault

import (
	"context"
	"testing"

	"github.com/aidapedia/gdk/config/secret"
)

func TestVault(t *testing.T) {
	ctx := context.Background()
	f := secret.NewSecretVault("your-address", string(secret.VaultEngineCubbyHole), "your-token", "your-path")

	var cfg struct {
		Auth struct {
			PrivateKey string
			PublicKey  string
		}
	}

	if err := f.GetSecret(ctx, &cfg); err != nil {
		t.Fatalf("Failed to get secret: %v", err)
	}
	if cfg.Auth.PrivateKey != "example-service" {
		t.Errorf("PrivateKey = %s; want %s", cfg.Auth.PrivateKey, "example-service")
	}
}
