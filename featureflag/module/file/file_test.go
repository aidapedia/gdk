package file

import (
	"context"
	"reflect"
	"testing"

	"github.com/aidapedia/gdk/featureflag/module"
	"github.com/bytedance/sonic"
)

func TestFeatureFlag_GetValue(t *testing.T) {
	type fields struct {
		configKeys Dir
	}
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "Value String",
			fields: fields{
				configKeys: Dir{
					Children: map[string]Dir{
						"featureflag": {
							Children: map[string]Dir{
								"test": {
									Children: map[string]Dir{},
									KVs: map[string]interface{}{
										"key1": "value1",
									},
								},
							},
						},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				key: "featureflag/test/key1",
			},
			want:    "value1",
			wantErr: false,
		},
		{
			name: "Not Found",
			fields: fields{
				configKeys: Dir{
					Children: map[string]Dir{
						"featureflag": {
							Children: map[string]Dir{
								"test": {
									Children: map[string]Dir{},
									KVs: map[string]interface{}{
										"key1": "value1",
									},
								},
							},
						},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				key: "featureflag/test/key2",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Value Integer",
			fields: fields{
				configKeys: Dir{
					Children: map[string]Dir{
						"featureflag": {
							Children: map[string]Dir{
								"test": {
									Children: map[string]Dir{},
									KVs: map[string]interface{}{
										"key1": 123,
									},
								},
							},
						},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				key: "featureflag/test/key1",
			},
			want:    123,
			wantErr: false,
		},
		{
			name: "Value Boolean",
			fields: fields{
				configKeys: Dir{
					Children: map[string]Dir{
						"featureflag": {
							Children: map[string]Dir{
								"test": {
									Children: map[string]Dir{},
									KVs: map[string]interface{}{
										"key1": true,
									},
								},
							},
						},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				key: "featureflag/test/key1",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Value JSON String",
			fields: fields{
				configKeys: Dir{
					Children: map[string]Dir{
						"featureflag": {
							Children: map[string]Dir{
								"test": {
									Children: map[string]Dir{},
									KVs: map[string]interface{}{
										"key1": "{\"key\":\"value\"}",
									},
								},
							},
						},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				key: "featureflag/test/key1",
			},
			want:    "{\"key\":\"value\"}",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &FeatureFlag{
				root: tt.fields.configKeys,
			}
			got, err := i.GetValue(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("FeatureFlag.GetValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FeatureFlag.GetValue() = %v, want %v", got, tt.want)
			}

			type TestStruct struct {
				Key string `json:"key"`
			}
			if tt.name == "Value JSON String" {
				var testStruct TestStruct
				err := sonic.Unmarshal([]byte(got.(string)), &testStruct)
				if err != nil {
					t.Errorf("FeatureFlag.GetValue() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if testStruct.Key != "value" {
					t.Errorf("FeatureFlag.GetValue() = %v, want %v", testStruct.Key, "value")
				}
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		filepath string
		prefix   string
	}
	tests := []struct {
		name string
		args args
		want module.Interface
	}{
		{
			name: "Normal Type",
			args: args{
				filepath: "test.json",
				prefix:   "/",
			},
			want: &FeatureFlag{
				root: Dir{
					Children: map[string]Dir{
						"featureflag": {
							Children: map[string]Dir{
								"test": {
									Children: map[string]Dir{},
									KVs: map[string]interface{}{
										"int": 2.0,
									},
								},
							},
						},
					},
					KVs: map[string]interface{}{
						"bool":        true,
						"string":      "string",
						"int":         1.0,
						"float":       1.0,
						"json_string": "{\"key\":\"json_string\"}",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.filepath, tt.args.prefix); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
