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

	// Optional fields, only present for tools responses
	ID       string                 `json:"id,omitempty"`
	Name     string                 `json:"name,omitempty"`
	Input    map[string]interface{} `json:"input,omitempty"`
	Thinking string                 `json:"thinking,omitempty"`
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
	Usage        MessageUsage          `json:"usage"`
}

type MessageUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

type MessageStreamResponse struct {
	Type  string             `json:"type"`
	Delta MessageStreamDelta `json:"delta"`
	Usage MessageStreamUsage `json:"usage"`
}

type MessageStreamDelta struct {
	Type         string `json:"type"`
	Text         string `json:"text"`
	StopReason   string `json:"stop_reason"`
	StopSequence string `json:"stop_sequence"`
	Thinking     string `json:"thinking"`
}

type MessageStreamUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}
