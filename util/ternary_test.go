package util

import (
	"testing"
)

func TestTernaryEqualString(t *testing.T) {
	type args struct {
		currentValue    string
		expressionValue string
		defaultVal      string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "equal",
			args: args{
				currentValue:    "a",
				expressionValue: "a",
				defaultVal:      "b",
			},
			want: "a",
		},
		{
			name: "not equal",
			args: args{
				currentValue:    "a",
				expressionValue: "c",
				defaultVal:      "b",
			},
			want: "b",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TernaryEqualString(tt.args.currentValue, tt.args.expressionValue, tt.args.defaultVal); got != tt.want {
				t.Errorf("TernaryEqualString() = %v, want %v", got, tt.want)
			}
		})
	}
}
