package anthropic

// CompletionRequestOption is a function type for CompletionRequest options
type CompletionRequestOption func(*CompletionRequest)

// MessageRequestOption is a function type for MessageRequest options
type MessageRequestOption func(*MessageRequest)

// Completion Request Options

func WithCompletionModel(model Model) CompletionRequestOption {
	return func(r *CompletionRequest) {
		r.Model = model
	}
}

func WithCompletionMaxTokens(maxTokens int) CompletionRequestOption {
	return func(r *CompletionRequest) {
		r.MaxTokensToSample = maxTokens
	}
}

func WithCompletionStream(stream bool) CompletionRequestOption {
	return func(r *CompletionRequest) {
		r.Stream = stream
	}
}

func WithCompletionStopSequences(stopSequences []string) CompletionRequestOption {
	return func(r *CompletionRequest) {
		r.StopSequences = stopSequences
	}
}

func WithCompletionTemperature(temperature float64) CompletionRequestOption {
	return func(r *CompletionRequest) {
		r.Temperature = temperature
	}
}

func WithCompletionTopK(topK int) CompletionRequestOption {
	return func(r *CompletionRequest) {
		r.TopK = topK
	}
}

func WithCompletionTopP(topP float64) CompletionRequestOption {
	return func(r *CompletionRequest) {
		r.TopP = topP
	}
}

// Message Request Options

func WithMessageModel(model Model) MessageRequestOption {
	return func(r *MessageRequest) {
		r.Model = model
	}
}

func WithMessages(messages []MessagePartRequest) MessageRequestOption {
	return func(r *MessageRequest) {
		r.Messages = messages
	}
}

func WithMessageMaxTokens(maxTokens int) MessageRequestOption {
	return func(r *MessageRequest) {
		r.MaxTokensToSample = maxTokens
	}
}

func WithSystemPrompt(systemPrompt string) MessageRequestOption {
	return func(r *MessageRequest) {
		r.SystemPrompt = systemPrompt
	}
}

func WithMetadata(metadata interface{}) MessageRequestOption {
	return func(r *MessageRequest) {
		r.Metadata = metadata
	}
}

func WithToolChoice(toolType, toolName string) MessageRequestOption {
	return func(r *MessageRequest) {
		r.ToolChoice = &ToolChoice{
			Type: toolType,
			Name: toolName,
		}
	}
}

func WithMessageStream(stream bool) MessageRequestOption {
	return func(r *MessageRequest) {
		r.Stream = stream
	}
}

func WithMessageStopSequences(stopSequences []string) MessageRequestOption {
	return func(r *MessageRequest) {
		r.StopSequences = stopSequences
	}
}

func WithMessageTemperature(temperature float64) MessageRequestOption {
	return func(r *MessageRequest) {
		r.Temperature = temperature
	}
}

func WithMessageTopK(topK int) MessageRequestOption {
	return func(r *MessageRequest) {
		r.TopK = topK
	}
}

func WithMessageTopP(topP float64) MessageRequestOption {
	return func(r *MessageRequest) {
		r.TopP = topP
	}
}
