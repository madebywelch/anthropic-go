package anthropic

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Complete sends a completion request to the API and returns a single completion.
func (c *Client) Complete(req *CompletionRequest) (*CompletionResponse, error) {
	if req.Stream {
		return nil, fmt.Errorf("cannot use Complete with a streaming request, use CompleteStream instead")
	}

	if !req.Model.IsCompleteCompatible() {
		return nil, fmt.Errorf("model %s is not compatible with the completion endpoint", req.Model)
	}

	return c.sendCompletionRequest(req)
}

func (c *Client) CompleteStream(req *CompletionRequest) (<-chan StreamResponse, <-chan error) {
	events := make(chan StreamResponse)

	// make this a buffered channel to allow for the error case below to return
	errCh := make(chan error, 1)

	if !req.Stream {
		errCh <- fmt.Errorf("cannot use CompleteStream with a non-streaming request, use Complete instead")
		return events, errCh
	}

	if !req.Model.IsCompleteCompatible() {
		errCh <- fmt.Errorf("model %s is not compatible with the completion endpoint", req.Model)
		return events, errCh
	}

	go c.handleStreaming(events, errCh, req)

	return events, errCh
}

// sendCompletionRequest sends a completion request to the API and returns a single completion.
func (c *Client) sendCompletionRequest(req *CompletionRequest) (*CompletionResponse, error) {
	// Marshal the request to JSON
	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshalling completion request: %w", err)
	}

	// Create the HTTP request
	requestURL := fmt.Sprintf("%s/v1/complete", c.baseURL)
	request, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(data))
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
	var completionResponse CompletionResponse
	err = json.NewDecoder(response.Body).Decode(&completionResponse)
	if err != nil {
		return nil, fmt.Errorf("error decoding completion response: %w", err)
	}

	return &completionResponse, nil
}

func (c *Client) handleStreaming(events chan StreamResponse, errCh chan error, req *CompletionRequest) {
	defer close(events)

	data, err := json.Marshal(req)
	if err != nil {
		errCh <- fmt.Errorf("error marshalling completion request: %w", err)
		return
	}

	request, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/complete", c.baseURL), bytes.NewBuffer(data))
	if err != nil {
		errCh <- fmt.Errorf("error creating new request: %w", err)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Api-Key", c.apiKey)
	request.Header.Set("Accept", "text/event-stream")

	response, err := c.doRequest(request)
	if err != nil {
		errCh <- fmt.Errorf("error sending completion request: %w", err)
		return
	}
	defer response.Body.Close()

	err = c.processSseStream(response.Body, events)
	if err != nil {
		errCh <- err
	}
}

func (c *Client) processSseStream(reader io.Reader, events chan StreamResponse) error {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "data:") {
			data := strings.TrimSpace(line[5:])
			var event StreamResponse
			err := json.Unmarshal([]byte(data), &event)
			if err != nil {
				return fmt.Errorf("error decoding event data: %w", err)
			}

			events <- event
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading from stream: %w", err)
	}

	return nil
}
