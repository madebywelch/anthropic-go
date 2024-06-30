package native

import (
	"net/http"

	"github.com/madebywelch/anthropic-go/v3/pkg/anthropic"
)

type Client struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
}

type Config struct {
	APIKey  string
	BaseURL string

	// Optional (defaults to http.DefaultClient)
	HTTPClient *http.Client
}

func MakeClient(cfg Config) (*Client, error) {
	if cfg.APIKey == "" {
		return nil, anthropic.ErrAnthropicApiKeyRequired
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.anthropic.com"
	}

	if cfg.HTTPClient == nil {
		cfg.HTTPClient = http.DefaultClient
	}

	return &Client{
		httpClient: cfg.HTTPClient,
		apiKey:     cfg.APIKey,
		baseURL:    cfg.BaseURL,
	}, nil
}
