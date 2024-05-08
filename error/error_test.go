package error

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		args []interface{}
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		{
			name: "Empty Error",
			args: args{
				args: []interface{}{},
			},
			want: &Error{},
		},
		{
			name: "All Set",
			args: args{
				args: []interface{}{
					Code(1),
					UserMessage("test"),
					map[string]interface{}{"param": "param-value"},
				},
			},
			want: &Error{
				code:        1,
				userMessage: "test",
				param: map[string]interface{}{
					"param": "param-value",
				},
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
