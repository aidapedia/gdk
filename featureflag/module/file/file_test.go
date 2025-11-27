package file

import (
	"context"
	"reflect"
	"testing"
)

func TestFeatureFlag_GetValue(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "key_bool",
			args: args{
				ctx: context.Background(),
				key: "key_bool",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "key_string",
			args: args{
				ctx: context.Background(),
				key: "key_string",
			},
			want:    "string",
			wantErr: false,
		},
		{
			name: "key_int",
			args: args{
				ctx: context.Background(),
				key: "key_int",
			},
			want:    float64(1),
			wantErr: false,
		},
		{
			name: "key_float",
			args: args{
				ctx: context.Background(),
				key: "key_float",
			},
			want:    1.0,
			wantErr: false,
		},
		{
			name: "key_json_string",
			args: args{
				ctx: context.Background(),
				key: "key_json_string",
			},
			want:    `{"key":"json_string"}`,
			wantErr: false,
		},
		{
			name: "folder",
			args: args{
				ctx: context.Background(),
				key: "folder/key",
			},
			want:    "value",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := New("test.json", "")
			got, err := i.GetValue(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("FeatureFlag.GetValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FeatureFlag.GetValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFeatureFlag_Watch(t *testing.T) {
	type fields struct {
		root    FolderItf
		address string
		prefix  string
	}
	tests := []struct {
		name    string
		fields  fields
		want    bool
		wantErr bool
	}{
		{
			name: "Changes",
			fields: fields{
				root: func() FolderItf {
					root, _ := parseFromFileJSON(nil, "test.json")
					root.Add(NewKey("changes", 2.0))
					return root
				}(),
				address: "test.json",
				prefix:  "/",
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &FeatureFlag{
				root:    tt.fields.root,
				address: tt.fields.address,
				prefix:  tt.fields.prefix,
			}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			got, err := i.Watch(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("FeatureFlag.Watch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			result := <-got
			if result != tt.want {
				t.Errorf("FeatureFlag.Watch() = %v, want %v", result, tt.want)
			}
			cancel()
		})
	}
}
