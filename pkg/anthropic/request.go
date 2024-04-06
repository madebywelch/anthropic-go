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
