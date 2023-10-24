package anthropic

// CompletionEvent represents the data related to a completion event.
type CompletionEvent struct {
	Completion string // The completion text from the model.
}

// ErrorEvent represents an error event that occurs during streaming.
type ErrorEvent struct {
	Error string // A string description of the error.
}

// PingEvent is an empty struct representing a ping event.
type PingEvent struct{}
