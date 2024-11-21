package anthropic

import (
	"errors"
	"fmt"
	"net/http"
)

type AnthropicErr struct {
	err         error
	errResponse string
}

func (e *AnthropicErr) Error() string {
	return fmt.Sprintf("%s", e.err)
}

func (e *AnthropicErr) Unwrap() error {
	return e.err
}

// ErrResponse is the upstream anthropic error response
func (e *AnthropicErr) ErrResponse() string {
	return e.errResponse
}

var (
	ErrAnthropicInvalidRequest = errors.New("invalid request: there was an issue with the format or content of your request")
	ErrAnthropicUnauthorized   = errors.New("unauthorized: there's an issue with your API key")
	ErrAnthropicForbidden      = errors.New("forbidden: your API key does not have permission to use the specified resource")
	ErrAnthropicRateLimit      = errors.New("your account has hit a rate limit")
	ErrAnthropicInternalServer = errors.New("an unexpected error has occurred internal to Anthropic's systems")
	ErrAnthropicApiKeyRequired = errors.New("apiKey is required")
	ErrAnthropicUnknown        = errors.New("unknown error occurred")
)

func NewAnthropicErr(err error, errResponse string) *AnthropicErr {
	return &AnthropicErr{
		err:         err,
		errResponse: errResponse,
	}
}

// mapHTTPStatusCodeToError maps an HTTP status code to an error.
func MapHTTPStatusCodeToError(code int, errResponse string) error {
	switch code {
	case http.StatusBadRequest:
		return &AnthropicErr{
			err:         ErrAnthropicInvalidRequest,
			errResponse: errResponse,
		}
	case http.StatusUnauthorized:
		return &AnthropicErr{
			err:         ErrAnthropicUnauthorized,
			errResponse: errResponse,
		}
	case http.StatusForbidden:
		return &AnthropicErr{
			err:         ErrAnthropicForbidden,
			errResponse: errResponse,
		}
	case http.StatusTooManyRequests:
		return &AnthropicErr{
			err:         ErrAnthropicRateLimit,
			errResponse: errResponse,
		}
	case http.StatusInternalServerError:
		return &AnthropicErr{
			err:         ErrAnthropicInternalServer,
			errResponse: errResponse,
		}
	default:
		return &AnthropicErr{
			err:         ErrAnthropicUnknown,
			errResponse: errResponse,
		}
	}
}
