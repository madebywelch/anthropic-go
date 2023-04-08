package anthropic

import (
	"bytes"
	"encoding/json"
	"fmt"
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

// Complete sends a completion request to the Anthropic API and returns the response.
func (c *Client) Complete(req *CompletionRequest) (*CompletionResponse, error) {
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

// GetPrompt returns a prompt string that can be used to complete a user question.
func GetPrompt(userQuestion string) string {
	const promptTemplate = `
Human: %s

Assistant:`
	return fmt.Sprintf(promptTemplate, strings.TrimSpace(userQuestion))
}
