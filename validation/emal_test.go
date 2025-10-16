package validation_test

import (
	"testing"

	"github.com/aidapedia/gdk/validation"
)

func TestIsEmail(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		email string
		want  bool
	}{
		{
			name:  "valid email",
			email: "test@example.com",
			want:  true,
		},
		{
			name:  "invalid email",
			email: "testexample",
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validation.IsEmail(tt.email)
			// TODO: update the condition below to compare got with tt.want.
			if got != tt.want {
				t.Errorf("IsEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}
