package utils

import (
	"testing"
)

func TestIsIPAllowed(t *testing.T) {
	tests := []struct {
		ip      string
		allowed bool
	}{
		{"160.79.104.1", true},
		{"160.79.103.1", false},
		{"2607:6bc0::1", true},
		{"2607:6bd0::1", false},
		{"invalid", false},
	}

	for _, test := range tests {
		result, _ := IsIPAllowed(test.ip)
		if result != test.allowed {
			t.Errorf("IsIPAllowed(%s) = %v; want %v", test.ip, result, test.allowed)
		}
	}
}
