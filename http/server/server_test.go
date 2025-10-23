package server_test

import (
	"testing"

	"github.com/aidapedia/gdk/http/server"
)

func TestNewWithDefaultConfig(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		opt     []server.Option
		want    *server.Server
		wantErr bool
	}{
		{
			name:    "success",
			opt:     []server.Option{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, gotErr := server.NewWithDefaultConfig("test-server", tt.opt...)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("NewWithDefaultConfig() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("NewWithDefaultConfig() succeeded unexpectedly")
			}
		})
	}
}
