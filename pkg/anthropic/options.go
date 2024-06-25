package anthropic

import "net/http"

type CompletionOption func(*CompletionRequest)
type MessageOption func(*MessageRequest)

type GenericOption[T any] func(*T)

func WithModel[T any](model Model) GenericOption[T] {
	return func(r *T) {
		switch v := any(r).(type) {
		case *CompletionRequest:
			v.Model = model
		case *MessageRequest:
			v.Model = model
		}
	}
}

func WithMessages[T MessageRequest](messages []MessagePartRequest) GenericOption[T] {
	return func(r *T) {
		switch v := any(r).(type) {
		case *MessageRequest:
			v.Messages = messages
		}
	}
}

func WithMaxTokens[T any](maxTokens int) GenericOption[T] {
	return func(r *T) {
		switch v := any(r).(type) {
		case *CompletionRequest:
			v.MaxTokensToSample = maxTokens
		case *MessageRequest:
			v.MaxTokensToSample = maxTokens
		}
	}
}

func WithSystemPrompt[T MessageRequest](systemPrompt string) GenericOption[T] {
	return func(r *T) {
		switch v := any(r).(type) {
		case *MessageRequest:
			v.SystemPrompt = systemPrompt
		}
	}
}

func WithMetadata[T MessageRequest](metadata interface{}) GenericOption[T] {
	return func(r *T) {
		switch v := any(r).(type) {
		case *MessageRequest:
			v.Metadata = metadata
		}
	}
}

func WithToolChoice[T MessageRequest](toolType, toolName string) GenericOption[T] {
	return func(r *T) {
		switch v := any(r).(type) {
		case *MessageRequest:
			v.ToolChoice = &ToolChoice{
				Type: toolType,
				Name: toolName,
			}
		}
	}
}

func WithStreaming[T any](stream bool) GenericOption[T] {
	return WithStream[T](stream)
}

func WithStream[T any](stream bool) GenericOption[T] {
	return func(r *T) {
		switch v := any(r).(type) {
		case *CompletionRequest:
			v.Stream = stream
		case *MessageRequest:
			v.Stream = stream
		}
	}
}

func WithStopSequences[T any](stopSequences []string) GenericOption[T] {
	return func(r *T) {
		switch v := any(r).(type) {
		case *CompletionRequest:
			v.StopSequences = stopSequences
		case *MessageRequest:
			v.StopSequences = stopSequences
		}
	}
}

func WithTemperature[T any](temperature float64) GenericOption[T] {
	return func(r *T) {
		switch v := any(r).(type) {
		case *CompletionRequest:
			v.Temperature = temperature
		case *MessageRequest:
			v.Temperature = temperature
		}
	}
}

func WithTopK[T any](topK int) GenericOption[T] {
	return func(r *T) {
		switch v := any(r).(type) {
		case *CompletionRequest:
			v.TopK = topK
		case *MessageRequest:
			v.TopK = topK
		}
	}
}

func WithTopP[T any](topP float64) GenericOption[T] {
	return func(r *T) {
		switch v := any(r).(type) {
		case *CompletionRequest:
			v.TopP = topP
		case *MessageRequest:
			v.TopP = topP
		}
	}
}

// WithHTTPClient sets a custom HTTP client for the Client.
func WithHTTPClient[T any](httpClient *http.Client) GenericOption[T] {
	return func(r *T) {
		if v, ok := any(r).(*Client); ok {
			v.httpClient = httpClient
		}
	}
}
