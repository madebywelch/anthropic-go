package anthropic

import (
	"errors"
	"net/http"
)

var (
	ErrAnthropicInvalidRequest = errors.New("invalid request: there was an issue with the format or content of your request")
	ErrAnthropicUnauthorized   = errors.New("unauthorized: there's an issue with your API key")
	ErrAnthropicForbidden      = errors.New("forbidden: your API key does not have permission to use the specified resource")
	ErrAnthropicRateLimit      = errors.New("your account has hit a rate limit")
	ErrAnthropicInternalServer = errors.New("an unexpected error has occurred internal to Anthropic's systems")

	ErrAnthropicApiKeyRequired = errors.New("apiKey is required")
)

// mapHTTPStatusCodeToError maps an HTTP status code to an error.
func mapHTTPStatusCodeToError(code int) error {
	switch code {
	case http.StatusBadRequest:
		return ErrAnthropicInvalidRequest
	case http.StatusUnauthorized:
		return ErrAnthropicUnauthorized
	case http.StatusForbidden:
		return ErrAnthropicForbidden
	case http.StatusTooManyRequests:
		return ErrAnthropicRateLimit
	case http.StatusInternalServerError:
		return ErrAnthropicInternalServer
	default:
		return errors.New("unknown error occurred")
	}
}
