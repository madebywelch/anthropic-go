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

// CompletionRequest is the request to the Anthropic API for a completion.
type CompletionRequest struct {
	Prompt            string   `json:"prompt"`
	Model             Model    `json:"model"`
	MaxTokensToSample int      `json:"max_tokens_to_sample"`
	StopSequences     []string `json:"stop_sequences,omitempty"` // optional
	Stream            bool     `json:"stream,omitempty"`         // optional
	Temperature       float64  `json:"temperature,omitempty"`    // optional
	TopK              int      `json:"top_k,omitempty"`          // optional
	TopP              float64  `json:"top_p,omitempty"`          // optional
}

// CompletionResponse is the response from the Anthropic API for a completion request.
type CompletionResponse struct {
	Completion string `json:"completion"`
	StopReason string `json:"stop_reason"`
	Stop       string `json:"stop"`
}

// StreamResponse is the response from the Anthropic API for a stream of completions.
type StreamResponse struct {
	Completion string `json:"completion"`
	StopReason string `json:"stop_reason"`
	Model      string `json:"model"`
	Stop       string `json:"stop"`
	LogID      string `json:"log_id"`
}

// StreamCallback is a function that handles a stream of completions.
type StreamCallback func(*CompletionResponse) error

// Complete sends a completion request to the API and returns a single completion or a stream of completions.
func (c *Client) Complete(req *CompletionRequest, callback StreamCallback) (*CompletionResponse, error) {
	if !req.Stream {
		response, err := c.sendCompletionRequest(req)
		if err != nil {
			return nil, err
		}
		return response, nil
	}

	return c.sendCompletionRequestStream(req, callback)
}

// sendCompletionRequest sends a completion request to the API and returns a single completion.
func (c *Client) sendCompletionRequest(req *CompletionRequest) (*CompletionResponse, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshalling completion request: %w", err)
	}

	request, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/complete", c.baseURL), bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("error creating new request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Api-Key", c.apiKey)

	response, err := c.doRequest(request)
	if err != nil {
		return nil, fmt.Errorf("error sending completion request: %w", err)
	}
	defer response.Body.Close()

	var completionResponse CompletionResponse
	err = json.NewDecoder(response.Body).Decode(&completionResponse)
	if err != nil {
		return nil, fmt.Errorf("error decoding completion response: %w", err)
	}

	return &completionResponse, nil
}

// sendCompletionRequestStream sends a completion request to the API and returns a stream of completions.
func (c *Client) sendCompletionRequestStream(req *CompletionRequest, callback StreamCallback) (*CompletionResponse, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshalling completion request: %w", err)
	}

	request, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/complete", c.baseURL), bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("error creating new request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Api-Key", c.apiKey)
	request.Header.Set("Accept", "text/event-stream")

	response, err := c.doRequest(request)
	if err != nil {
		return nil, fmt.Errorf("error sending completion request: %w", err)
	}
	defer response.Body.Close()

	return c.processSseStream(response.Body, callback)
}

// processSseStream reads the SSE stream from the API and calls the callback for each completion.
func (c *Client) processSseStream(reader io.Reader, callback StreamCallback) (*CompletionResponse, error) {
	scanner := bufio.NewScanner(reader)
	var completionResponse CompletionResponse

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data:") {
			data := strings.TrimSpace(line[5:])

			var event StreamResponse
			err := json.Unmarshal([]byte(data), &event)
			if err != nil {
				return nil, fmt.Errorf("error decoding event data: %w", err)
			}

			completionResponse.Completion += event.Completion
			completionResponse.StopReason = event.StopReason
			completionResponse.Stop = event.Stop

			if err := callback(&completionResponse); err != nil {
				return nil, err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading from stream: %w", err)
	}

	return &completionResponse, nil
}
