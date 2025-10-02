package file

import (
	"context"
	"testing"

	"github.com/aidapedia/gdk/config/secret"
)

func TestFile_GetSecret(t *testing.T) {
	ctx := context.Background()
	f := secret.NewSecretFile("config/secret.json")

	var cfg struct {
		ServiceName string
	}

	if err := f.GetSecret(ctx, &cfg); err != nil {
		t.Fatalf("Failed to get secret: %v", err)
	}
	if cfg.ServiceName != "example-service" {
		t.Errorf("ServiceName = %s; want %s", cfg.ServiceName, "example-service")
	}
}
