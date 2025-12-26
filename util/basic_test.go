package util

import (
	"encoding/json"
	"testing"
	"time"
)

func TestToStr(t *testing.T) {
	tests := []struct {
		name string
		arg  interface{}
		want string
	}{
		{"nil", nil, ""},
		{"string", "abc", "abc"},
		{"int", 123, "123"},
		{"int64", int64(123), "123"},
		{"bool", true, "true"},
		{"float", 1.23, "1.23"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToStr(tt.arg); got != tt.want {
				t.Errorf("ToStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToBool(t *testing.T) {
	tests := []struct {
		name string
		arg  interface{}
		want bool
	}{
		{"string true", "true", true},
		{"string false", "false", false},
		{"int 1", 1, true},
		{"int 0", 0, false},
		{"bool true", true, true},
		{"nil", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToBool(tt.arg); got != tt.want {
				t.Errorf("ToBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToInt(t *testing.T) {
	tests := []struct {
		name string
		arg  interface{}
		want int
	}{
		{"string 123", "123", 123},
		{"string invalid", "abc", 0},
		{"int 123", 123, 123},
		{"float 1.23", 1.23, 1},
		{"bool true", true, 1},
		{"nil", nil, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToInt(tt.arg); got != tt.want {
				t.Errorf("ToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToInt64(t *testing.T) {
	tests := []struct {
		name string
		arg  interface{}
		want int64
	}{
		{"string 123", "123", 123},
		{"int 123", 123, 123},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToInt64(tt.arg); got != tt.want {
				t.Errorf("ToInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToFloat64(t *testing.T) {
	tests := []struct {
		name string
		arg  interface{}
		want float64
	}{
		{"string 1.23", "1.23", 1.23},
		{"int 1", 1, 1.0},
		{"json", json.RawMessage(`1.23`), 1.23},
		{"nil", nil, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToFloat64(tt.arg); got != tt.want {
				t.Errorf("ToFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToTime(t *testing.T) {
	now := time.Now()
	nowStr := now.Format(time.RFC3339)
	parsed, _ := time.Parse(time.RFC3339, nowStr)

	tests := []struct {
		name string
		arg  interface{}
		want time.Time
	}{
		{"time", now, now},
		{"string", nowStr, parsed},
		{"invalid", "invalid", time.Time{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToTime(tt.arg); !got.Equal(tt.want) {
				t.Errorf("ToTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
