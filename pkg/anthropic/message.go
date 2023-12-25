package anthropic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) Message(req *MessageRequest) (*MessageResponse, error) {
	if req.Stream {
		return nil, fmt.Errorf("cannot use Message with streaming enabled, use MessageStream instead (not yet supported)")
	}

	return c.sendMessageRequest(req)
}

// MessageStream (NOT YET SUPPORTED) returns a channel of StreamResponse objects and a channel of errors.
func (c *Client) MessageStream(req *MessageRequest) (<-chan StreamResponse, <-chan error) {
	return nil, nil
}

func (c *Client) sendMessageRequest(req *MessageRequest) (*MessageResponse, error) {
	// Marshal the request to JSON
	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshalling completion request: %w", err)
	}

	// Create the HTTP request
	requestURL := fmt.Sprintf("%s/v1/messages", c.baseURL)
	request, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("error creating new request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Api-Key", c.apiKey)
	request.Header.Set("anthropic-beta", AnthropicAPIMessagesBeta)

	// Use the DoRequest method to send the HTTP request
	response, err := c.doRequest(request)
	if err != nil {
		return nil, fmt.Errorf("error sending completion request: %w", err)
	}
	defer response.Body.Close()

	// Decode the response body to a MessageResponse object
	var messageResponse MessageResponse
	err = json.NewDecoder(response.Body).Decode(&messageResponse)
	if err != nil {
		return nil, fmt.Errorf("error decoding message response: %w", err)
	}

	return &messageResponse, nil
}
