package anthropic

// CompletionRequest is the request to the Anthropic API for a completion.
type CompletionRequest struct {
	Prompt            string   `json:"prompt"`
	Model             Model    `json:"model,omitempty"`
	MaxTokensToSample int      `json:"max_tokens_to_sample"`
	StopSequences     []string `json:"stop_sequences,omitempty"` // optional
	Stream            bool     `json:"stream,omitempty"`         // optional
	Temperature       float64  `json:"temperature,omitempty"`    // optional
	TopK              int      `json:"top_k,omitempty"`          // optional
	TopP              float64  `json:"top_p,omitempty"`          // optional
}

// NewCompletionRequest creates a new CompletionRequest with the given prompt and options.
//
// prompt: the prompt for the completion request.
// options: optional GenericOptions to customize the completion request.
// Returns a pointer to the newly created CompletionRequest.
func NewCompletionRequest(prompt string, options ...GenericOption[CompletionRequest]) *CompletionRequest {
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

// ContentBlock interface to allow for both TextContentBlock and ImageContentBlock
type ContentBlock interface {
	// This method exists solely to enforce compile-time checking of the types that implement this interface.
	isContentBlock()
}

// TextContentBlock represents a block of text content.
type TextContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func (t TextContentBlock) isContentBlock() {}

// ImageSource represents the source of an image, supporting base64 encoding for now.
type ImageSource struct {
	Type      string    `json:"type"`
	MediaType MediaType `json:"media_type"`
	Data      string    `json:"data"`
}

// ImageContentBlock represents a block of image content.
type ImageContentBlock struct {
	Type   string      `json:"type"`
	Source ImageSource `json:"source"`
}

func (i ImageContentBlock) isContentBlock() {}

// MessagePartRequest is updated to support both text and image content blocks.
type MessagePartRequest struct {
	Role    string         `json:"role"`
	Content []ContentBlock `json:"content"`
}

// Helper functions to create text and image content blocks easily
func NewTextContentBlock(text string) ContentBlock {
	return TextContentBlock{
		Type: "text",
		Text: text,
	}
}

type MediaType string

const (
	MediaTypeJPEG MediaType = "image/jpeg"
	MediaTypePNG  MediaType = "image/png"
	MediaTypeGIF  MediaType = "image/gif"
	MediaTypeWEBP MediaType = "image/webp"
)

// NewImageContentBlock creates a new image content block with the given media type and base64 data.
func NewImageContentBlock(mediaType MediaType, base64Data string) ContentBlock {
	return ImageContentBlock{
		Type: "image",
		Source: ImageSource{
			Type:      "base64",
			MediaType: mediaType,
			Data:      base64Data,
		},
	}
}

// ToolResultContentBlock represents a block of tool result content.
type ToolResultContentBlock struct {
	Type      string      `json:"type"`
	ToolUseID string      `json:"tool_use_id"`
	Content   interface{} `json:"content"`
	IsError   bool        `json:"is_error,omitempty"`
}

func (t ToolResultContentBlock) isContentBlock() {}

// NewToolResultContentBlock creates a new tool result content block with the given parameters.
func NewToolResultContentBlock(toolUseID string, content interface{}, isError bool) ContentBlock {
	return ToolResultContentBlock{
		Type:      "tool_result",
		ToolUseID: toolUseID,
		Content:   content,
		IsError:   isError,
	}
}

// ToolChoice specifies the tool preferences for a message request.
type ToolChoice struct {
	Type string `json:"type"` // Type of tool choice: "tool", "any", or "auto".
	Name string `json:"name"` // Name of the tool to be used (if type is "tool").
}

// MessageRequest is the request to the Anthropic API for a message request.
type MessageRequest struct {
	Model             Model                `json:"model"`
	Tools             []Tool               `json:"tools,omitempty"`
	Messages          []MessagePartRequest `json:"messages"`
	MaxTokensToSample int                  `json:"max_tokens"`
	SystemPrompt      string               `json:"system,omitempty"`         // optional
	Metadata          interface{}          `json:"metadata,omitempty"`       // optional
	StopSequences     []string             `json:"stop_sequences,omitempty"` // optional
	Stream            bool                 `json:"stream,omitempty"`         // optional
	Temperature       float64              `json:"temperature,omitempty"`    // optional
	ToolChoice        *ToolChoice          `json:"tool_choice,omitempty"`    // optional
	TopK              int                  `json:"top_k,omitempty"`          // optional
	TopP              float64              `json:"top_p,omitempty"`          // optional
}

type Property struct {
	Type        string   `json:"type"`
	Enum        []string `json:"enum,omitempty"`
	Description string   `json:"description"`
}

type InputSchema struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
	Required   []string            `json:"required"`
}

type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema InputSchema `json:"input_schema"`
}

// CountImageContent counts the number of ImageContentBlock in the MessageRequest.
//
// No parameters.
// Returns an integer representing the count.
func (m *MessageRequest) CountImageContent() int {
	count := 0
	for _, message := range m.Messages {
		for _, block := range message.Content {
			if _, ok := block.(ImageContentBlock); ok {
				count++
			}
		}
	}
	return count
}

// ContainsImageContent checks if the MessageRequest contains any ImageContentBlock.
//
// No parameters.
// Returns a boolean value.
func (m *MessageRequest) ContainsImageContent() bool {
	for _, message := range m.Messages {
		for _, block := range message.Content {
			if _, ok := block.(ImageContentBlock); ok {
				return true
			}
		}
	}
	return false
}

// NewMessageRequest creates a new MessageRequest with the provided messages and options.
// It takes in a slice of MessagePartRequests and optional GenericOptions and returns a pointer to a MessageRequest.
func NewMessageRequest(messages []MessagePartRequest, options ...GenericOption[MessageRequest]) *MessageRequest {
	request := &MessageRequest{
		Messages: messages,
		// defauts, can be overridden
		Model:             ClaudeV2,
		MaxTokensToSample: 25,
	}
	for _, option := range options {
		option(request)
	}
	return request
}
