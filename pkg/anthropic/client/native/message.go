package native

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/madebywelch/anthropic-go/v3/pkg/anthropic"
)

func (c *Client) Message(ctx context.Context, req *anthropic.MessageRequest) (*anthropic.MessageResponse, error) {
	err := anthropic.ValidateMessageRequest(req)
	if err != nil {
		return nil, err
	}

	return c.sendMessageRequest(ctx, req)
}

func (c *Client) sendMessageRequest(
	ctx context.Context,
	req *anthropic.MessageRequest,
) (*anthropic.MessageResponse, error) {
	// Marshal the request to JSON
	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshalling message request: %w", err)
	}

	// Create the HTTP request
	requestURL := fmt.Sprintf("%s/v1/messages", c.baseURL)
	request, err := http.NewRequestWithContext(ctx, "POST", requestURL, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("error creating new request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Api-Key", c.apiKey)
	if len(req.Beta) > 0 {
		request.Header.Set("anthropic-beta", req.Beta)
	}

	// Use the doRequest method to send the HTTP request
	response, err := c.doRequest(request)
	if err != nil {
		return nil, fmt.Errorf("error sending message request: %w", err)
	}
	defer response.Body.Close()

	// Decode the response body to a MessageResponse object
	messageResponse := &anthropic.MessageResponse{}
	err = json.NewDecoder(response.Body).Decode(messageResponse)
	if err != nil {
		return nil, fmt.Errorf("error decoding message response: %w", err)
	}

	return messageResponse, nil
}
