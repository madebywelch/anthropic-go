package anthropic

import (
	"errors"
	"net/http"
	"testing"
)

func TestMapHTTPStatusCodeToError(t *testing.T) {
	tests := []struct {
		code     int
		expected error
	}{
		{http.StatusBadRequest, ErrAnthropicInvalidRequest},
		{http.StatusUnauthorized, ErrAnthropicUnauthorized},
		{http.StatusForbidden, ErrAnthropicForbidden},
		{http.StatusTooManyRequests, ErrAnthropicRateLimit},
		{http.StatusInternalServerError, ErrAnthropicInternalServer},
		{http.StatusNotFound, errors.New("unknown error occurred")},
	}

	for _, test := range tests {
		err := MapHTTPStatusCodeToError(test.code)
		if err.Error() != test.expected.Error() {
			t.Errorf("Expected error '%s', got '%s'", test.expected.Error(), err.Error())
		}
	}
}
