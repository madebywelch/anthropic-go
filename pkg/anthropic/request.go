package anthropic

import "github.com/invopop/jsonschema"

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
func NewCompletionRequest(prompt string, options ...CompletionRequestOption) *CompletionRequest {
	request := &CompletionRequest{
		Prompt:            prompt,
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

func (t ToolUseContentBlock) isContentBlock()    {}
func (t ToolResultContentBlock) isContentBlock() {}

// ToolUseContentBlock represents a block of tool use content.
type ToolUseContentBlock struct {
	Type  string      `json:"type"`
	ID    string      `json:"id"`
	Name  string      `json:"name"`
	Input interface{} `json:"input"`
}

func NewToolUseContentBlock(id string, name string, input interface{}) ContentBlock {
	return ToolUseContentBlock{
		Type: "tool_use",
		ID:   id,
		Name: name,
		Input: map[string]interface{}{
			"type":  name,
			"input": input,
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
	Type string `json:"type"`           // Type of tool choice: "tool", "any", or "auto".
	Name string `json:"name,omitempty"` // Name of the tool to be used (if type is "tool").
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

type Tool struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	InputSchema *jsonschema.Schema `json:"input_schema"`
}

func GenerateInputSchema(input interface{}) *jsonschema.Schema {
	return (&jsonschema.Reflector{ExpandedStruct: true}).Reflect(input)
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

// NewMessageRequest creates a new MessageRequest with the provided options.
func NewMessageRequest(options ...MessageRequestOption) *MessageRequest {
	request := &MessageRequest{
		Model:             ClaudeV2,
		MaxTokensToSample: 25,
	}
	for _, option := range options {
		option(request)
	}
	return request
}

// AddMessage adds a message to the MessageRequest.
func (r *MessageRequest) AddMessage(role string, content ...ContentBlock) *MessageRequest {
	r.Messages = append(r.Messages, MessagePartRequest{
		Role:    role,
		Content: content,
	})
	return r
}

// AddUserMessage adds a user message to the MessageRequest.
func (r *MessageRequest) AddUserMessage(content ...ContentBlock) *MessageRequest {
	return r.AddMessage("user", content...)
}

// AddAssistantMessage adds an assistant message to the MessageRequest.
func (r *MessageRequest) AddAssistantMessage(content ...ContentBlock) *MessageRequest {
	return r.AddMessage("assistant", content...)
}
