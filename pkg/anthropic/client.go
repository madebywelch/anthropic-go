package anthropic

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

// Client represents the Anthropic API client and its configuration.
type Client struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
	maxRetries int
	retryDelay time.Duration
}

// ClientOptions is a variadic function that configures the client.
type ClientOptions func(*Client)

// WithMaxRetries sets the maximum number of retries for a request.
func WithMaxRetries(maxRetries int) ClientOptions {
	return func(c *Client) {
		c.maxRetries = maxRetries
	}
}

// WithRetryDelay sets the delay between retries for a request.
func WithRetryDelay(retryDelay time.Duration) ClientOptions {
	return func(c *Client) {
		c.retryDelay = retryDelay
	}
}

// NewClient initializes a new Anthropic API client with the required headers.
func NewClient(apiKey string, options ...ClientOptions) (*Client, error) {
	if apiKey == "" {
		return nil, errors.New("apiKey is required")
	}

	client := &Client{
		httpClient: &http.Client{},
		apiKey:     apiKey,
		baseURL:    "https://api.anthropic.com",
		maxRetries: 0,               // Default value for maxRetries
		retryDelay: 0 * time.Second, // Default value for retryDelay
	}

	for _, option := range options {
		option(client)
	}

	return client, nil
}

// doRequest is a wrapper function that handles retries and delays for all API calls.
func (c *Client) doRequest(request *http.Request) (*http.Response, error) {
	var (
		response *http.Response
		err      error
	)

	retries := c.maxRetries
	if retries < 0 {
		retries = 0
	}

	delay := c.retryDelay
	if delay < 0 {
		delay = 0
	}

	for i := 0; i <= retries; i++ {
		response, err = c.httpClient.Do(request)
		if err == nil {
			if response.StatusCode == http.StatusOK {
				break
			} else {
				err = fmt.Errorf("unexpected status code: %d", response.StatusCode)
			}
		}

		if i < retries {
			time.Sleep(delay)
		}
	}

	if err != nil {
		return nil, err
	}

	return response, nil
}
