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

func TestManager_GetConfig(t *testing.T) {
	type fields struct {
		store      interface{}
		secretType SecretType
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
			name: "GetConfig with store set",
			fields: fields{
				store: &testConfig{
					Environment: "dev",
				},
				secretType: SecretTypeFile,
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
			wantTarget: &testConfig{
				Environment: "dev",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{
				store:      tt.fields.store,
				secretType: tt.fields.secretType,
			}
			store, err := m.GetConfig(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.GetConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
			cf := store.(*testConfig)
			if !reflect.DeepEqual(cf, tt.wantTarget) {
				t.Errorf("Manager.GetConfig() store = %v, wantTarget %v", store, tt.wantTarget)
			}
		})
	}
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
			cf, err := m.GetConfig(tt.args.ctx)
			if err != nil {
				t.Errorf("Manager.GetConfig() error = %v", err)
			}
			cf = cf.(*testConfig)
			if !reflect.DeepEqual(cf, tt.wantTarget) {
				t.Errorf("Manager.GetConfig() store = %v, wantTarget %v", cf, tt.wantTarget)
			}
		})
	}
}
