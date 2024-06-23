package anthropic

import (
	"net/http"
)

// Client represents the Anthropic API client and its configuration.
type Client struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
}

// NewClient initializes a new Anthropic API client with the required headers.
func NewClient(apiKey string, options ...GenericOption[Client]) (*Client, error) {
	if apiKey == "" {
		return nil, ErrAnthropicApiKeyRequired
	}

	client := &Client{
		httpClient: &http.Client{},
		apiKey:     apiKey,
		baseURL:    "https://api.anthropic.com",
	}
	for _, opt := range options {
		opt(client)
	}

	return client, nil
}
