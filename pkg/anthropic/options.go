package anthropic

type CompletionOption func(*CompletionRequest)

func WithModel(model Model) CompletionOption {
	return func(r *CompletionRequest) {
		r.Model = model
	}
}

func WithMaxTokens(maxTokens int) CompletionOption {
	return func(r *CompletionRequest) {
		r.MaxTokensToSample = maxTokens
	}
}

func WithStopSequences(stopSequences []string) CompletionOption {
	return func(r *CompletionRequest) {
		r.StopSequences = stopSequences
	}
}

func WithStreaming(stream bool) CompletionOption {
	return func(r *CompletionRequest) {
		r.Stream = stream
	}
}

func WithTemperature(temperature float64) CompletionOption {
	return func(r *CompletionRequest) {
		r.Temperature = temperature
	}
}

func WithTopK(topK int) CompletionOption {
	return func(r *CompletionRequest) {
		r.TopK = topK
	}
}

func WithTopP(topP float64) CompletionOption {
	return func(r *CompletionRequest) {
		r.TopP = topP
	}
}
