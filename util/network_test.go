package util_test

import (
	"testing"

	"github.com/aidapedia/gdk/util"
)

func TestCheckSubnet(t *testing.T) {
	tests := []struct {
		name   string
		ip     string
		subnet string
		want   bool
	}{
		{
			name:   "ip in subnet",
			ip:     "192.168.1.1",
			subnet: "192.168.1.0/24",
			want:   true,
		},
		{
			name:   "ip not in subnet",
			ip:     "192.168.2.1",
			subnet: "192.168.1.0/24",
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := util.CheckSubnet(tt.ip, tt.subnet)
			if got != tt.want {
				t.Errorf("CheckSubnet() = %v, want %v", got, tt.want)
			}
		})
	}
}
