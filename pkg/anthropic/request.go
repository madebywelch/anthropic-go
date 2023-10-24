package anthropic

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

func NewCompletionRequest(prompt string, options ...CompletionOption) *CompletionRequest {
	request := &CompletionRequest{
		Prompt: prompt,
		// defauts, can be overridden
		Model:             ClaudeV2,
		MaxTokensToSample: 25,
	}
	for _, option := range options {
		option(request)
	}
	return request
}
