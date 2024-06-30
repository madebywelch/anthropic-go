package native

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/madebywelch/anthropic-go/v3/pkg/anthropic"
)

func (c *Client) Complete(ctx context.Context, req *anthropic.CompletionRequest) (*anthropic.CompletionResponse, error) {
	err := anthropic.ValidateCompleteRequest(req)
	if err != nil {
		return nil, err
	}

	return c.sendCompleteRequest(ctx, req)
}

func (c *Client) sendCompleteRequest(ctx context.Context, req *anthropic.CompletionRequest) (*anthropic.CompletionResponse, error) {
	// Marshal the request to JSON
	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshalling completion request: %w", err)
	}

	// Create the HTTP request
	requestURL := fmt.Sprintf("%s/v1/complete", c.baseURL)
	request, err := http.NewRequestWithContext(ctx, "POST", requestURL, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("error creating new request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Api-Key", c.apiKey)

	// Use the DoRequest method to send the HTTP request
	response, err := c.doRequest(request)
	if err != nil {
		return nil, fmt.Errorf("error sending completion request: %w", err)
	}
	defer response.Body.Close()

	// Decode the response body to a CompletionResponse object
	completionResponse := &anthropic.CompletionResponse{}
	err = json.NewDecoder(response.Body).Decode(&completionResponse)
	if err != nil {
		return nil, fmt.Errorf("error decoding completion response: %w", err)
	}

	return completionResponse, nil
}
