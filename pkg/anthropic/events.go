package anthropic

// Common types for different events
type MessageEventType string

// Define a separate type for completion events
type CompletionEventType string

const (
	// Constants for message event types
	MessageEventTypeMessageStart      MessageEventType = "message_start"
	MessageEventTypeContentBlockStart MessageEventType = "content_block_start"
	MessageEventTypePing              MessageEventType = "ping"
	MessageEventTypeContentBlockDelta MessageEventType = "content_block_delta"
	MessageEventTypeContentBlockStop  MessageEventType = "content_block_stop"
	MessageEventTypeMessageDelta      MessageEventType = "message_delta"
	MessageEventTypeMessageStop       MessageEventType = "message_stop"
	MessageEventTypeError             MessageEventType = "error"

	// Constants for completion event types
	CompletionEventTypeCompletion CompletionEventType = "completion"
	CompletionEventTypePing       CompletionEventType = "ping"
)
