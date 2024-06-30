package utils

import (
	"testing"
)

func TestIsRegionSupported(t *testing.T) {
	regionTests := map[string]bool{
		"United States of America": true,
		"Kamino":                   false,
		"Germany":                  true,
		"":                         false,
	}

	for region, expected := range regionTests {
		actual := IsRegionSupported(region)
		if actual != expected {
			t.Errorf("IsRegionSupported(%q) = %v; expected %v", region, actual, expected)
		}
	}
}
