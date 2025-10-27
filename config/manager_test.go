package config

import (
	"context"
	"os"
	"reflect"
	"testing"
)

type testConfig struct {
	Environment string
	AppEnv      string
}

func TestManager_SetConfig(t *testing.T) {
	type fields struct {
		store    interface{}
		fileName []string
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		wantTarget *testConfig
	}{
		{
			name: "SetConfig with store set",
			fields: fields{
				fileName: []string{"config", "app"},
				store:    &testConfig{},
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
			wantTarget: &testConfig{
				Environment: "dev",
				AppEnv:      "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := os.Setenv("CONFIG_FILE_PATH", "test")
			if err != nil {
				t.Errorf("os.Setenv() error = %v", err)
			}

			m := New(Option{
				TargetStore: tt.fields.store,
				FileName:    tt.fields.fileName,
				ConfigKey:   "Config",
			})
			err = m.SetConfig(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.SetConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.fields.store, tt.wantTarget) {
				t.Errorf("Manager.SetConfig() store = %v, wantTarget %v", tt.fields.store, tt.wantTarget)
			}
		})
	}
}
