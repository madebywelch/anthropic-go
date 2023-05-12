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
	Delta      string `json:"delta,omitempty"`
	StopReason string `json:"stop_reason"`
	Stop       string `json:"stop"`
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
	var dataBuffer bytes.Buffer
	var prev string

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, ":") {
			continue
		}

		if line == "data: [DONE]" {
			break
		}

		if strings.HasPrefix(line, "data:") {
			data := strings.TrimSpace(line[5:])
			dataBuffer.WriteString(data)
		} else if line == "" {
			if dataBuffer.Len() > 0 {
				var completionResponse CompletionResponse
				err := json.Unmarshal(dataBuffer.Bytes(), &completionResponse)
				dataBuffer.Reset()

				if prev != "" {
					arr := strings.SplitAfter(completionResponse.Completion, prev)

					if len(arr) > 1 {
						completionResponse.Delta = arr[1]
					} else {
						return nil, fmt.Errorf("could not compute delta")
					}
				} else {
					completionResponse.Delta = completionResponse.Completion
				}

				prev = completionResponse.Completion

				if err != nil {
					return nil, fmt.Errorf("error decoding completion response: %w", err)
				}
				if err := callback(&completionResponse); err != nil {
					return nil, err
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading from stream: %w", err)
	}

	return nil, nil
}

// GetPrompt returns a prompt string that can be used to complete a user question.
func GetPrompt(userQuestion string) string {
	const promptTemplate = `
Human: %s

Assistant:`
	return fmt.Sprintf(promptTemplate, strings.TrimSpace(userQuestion))
}
