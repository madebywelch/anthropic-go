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
func NewClient(apiKey string) (*Client, error) {
	if apiKey == "" {
		return nil, ErrAnthropicApiKeyRequired
	}

	client := &Client{
		httpClient: &http.Client{},
		apiKey:     apiKey,
		baseURL:    "https://api.anthropic.com",
	}

	return client, nil
}

// doRequest is a wrapper function that handles retries and delays for all API calls.
func (c *Client) doRequest(request *http.Request) (*http.Response, error) {

	request.Header.Add("anthropic-version", "2023-01-01") // This SDK is not compatible with 2023-06-01

	var (
		response *http.Response
		err      error
	)

	response, err = c.httpClient.Do(request)
	if err == nil {
		if response.StatusCode == http.StatusOK {
			return response, nil
		} else {
			err = mapHTTPStatusCodeToError(response.StatusCode)
		}
	}

	if err != nil {
		return nil, err
	}

	return response, nil
}
