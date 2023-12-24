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

// MessageResponse is a subset of the response from the Anthropic API for a message response.
type MessagePartResponse struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// MessageResponse is the response from the Anthropic API for a message response.
type MessageResponse struct {
	ID           string                `json:"id"`
	Type         string                `json:"type"`
	Model        string                `json:"model"`
	Role         string                `json:"role"`
	Content      []MessagePartResponse `json:"content"`
	StopReason   string                `json:"stop_reason"`
	Stop         string                `json:"stop"`
	StopSequence string                `json:"stop_sequence"`
}
