package mask

import (
	"reflect"
	"testing"

	masker "github.com/ggwhite/go-masker/v2"
)

func TestMask_MaskMap(t *testing.T) {
	type TestStruct struct {
		Email string `json:"email" mask:"email"`
	}

	type fields struct {
		MaskerMarshaler *masker.MaskerMarshaler
	}
	type args struct {
		val map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "Normal Type",
			fields: fields{
				MaskerMarshaler: masker.NewMaskerMarshaler(),
			},
			args: args{
				val: map[string]interface{}{
					"email": "abc@gmail.com",
				},
			},
			want: map[string]interface{}{
				"email": "abc****@gmail.com",
			},
			wantErr: false,
		},
		{
			name: "Map Type",
			fields: fields{
				MaskerMarshaler: masker.NewMaskerMarshaler(),
			},
			args: args{
				val: map[string]interface{}{
					"email": map[string]interface{}{
						"email": "abc@gmail.com",
					},
				},
			},
			want: map[string]interface{}{
				"email": map[string]interface{}{
					"email": "abc****@gmail.com",
				},
			},
			wantErr: false,
		},
		{
			name: "Struct Type",
			fields: fields{
				MaskerMarshaler: masker.NewMaskerMarshaler(),
			},
			args: args{
				val: map[string]interface{}{
					"email": TestStruct{
						Email: "abc@gmail.com",
					},
				},
			},
			want: map[string]interface{}{
				"email": &TestStruct{
					Email: "abc****@gmail.com",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Mask{
				MaskerMarshaler: tt.fields.MaskerMarshaler,
			}
			got, err := m.MaskMap(tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("Mask.MaskMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Mask.MaskMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
