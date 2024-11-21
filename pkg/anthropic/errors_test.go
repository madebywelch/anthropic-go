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
		fullErr  string
	}{
		{http.StatusBadRequest, ErrAnthropicInvalidRequest, "test err body 1"},
		{http.StatusUnauthorized, ErrAnthropicUnauthorized, "test err body 2"},
		{http.StatusForbidden, ErrAnthropicForbidden, "test err body 3"},
		{http.StatusTooManyRequests, ErrAnthropicRateLimit, "test err body 4"},
		{http.StatusInternalServerError, ErrAnthropicInternalServer, "test err body 5"},
		{http.StatusNotFound, errors.New("unknown error occurred"), "test err body 6"},
	}

	for _, test := range tests {
		err := MapHTTPStatusCodeToError(test.code, test.fullErr)
		if err.Error() != test.expected.Error() {
			t.Errorf("Expected error '%s', got '%s'", test.expected.Error(), err.Error())
		}

		var actualErr *AnthropicErr
		if errors.As(err, &actualErr) {
			if actualErr.ErrResponse() != test.fullErr {
				t.Errorf("full error response doesn't match. expected: %s | actual: %s", test.fullErr, actualErr.ErrResponse())
			}
		} else {
			t.Error("error is not of type `anthropic.AnthropicErr`")
		}
	}
}
