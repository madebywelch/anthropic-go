package anthropic

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
