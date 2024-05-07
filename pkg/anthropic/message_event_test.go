package anthropic

import (
	"reflect"
	"testing"
)

func TestParseMessageEvent(t *testing.T) {
	events := []struct {
		eventType MessageEventType
		event     string
		expected  *MessageStreamResponse
		expErrStr string
	}{
		{
			eventType: MessageEventTypeMessageStart,
			event: `{
				"type": "message_start",
				"message": {
					"id": "123",
					"type": "text",
					"role": "user",
					"content": ["Hello, world!"],
					"model": "claude-v2_1",
					"stop_reason": "",
					"stop_sequence": "",
					"usage": {
						"input_tokens": 10,
						"output_tokens": 20
					}
				}
			}`,
			expected: &MessageStreamResponse{
				Type: "message_start",
				Usage: MessageStreamUsage{
					InputTokens:  10,
					OutputTokens: 20,
				},
			},
		},
		{
			eventType: MessageEventTypeContentBlockStart,
			event: `{
				"type": "content_block_start",
				"index": 1,
				"content_block": {
					"type": "text",
					"text": "This is a content block"
				}
			}`,
			expected: &MessageStreamResponse{
				Type: "content_block_start",
			},
		},
		{
			eventType: MessageEventTypePing,
			event: `{
				"type": "ping"
			}`,
			expected: &MessageStreamResponse{
				Type: "ping",
			},
		},
		{
			eventType: MessageEventTypeContentBlockDelta,
			event: `{
				"type": "content_block_delta",
				"delta": {
					"type": "text",
					"text": "This is a content block delta"
				}
			}`,
			expected: &MessageStreamResponse{
				Type: "content_block_delta",
				Delta: MessageStreamDelta{
					Type:         "text",
					Text:         "This is a content block delta",
					StopReason:   "",
					StopSequence: "",
				},
			},
		},
		{
			eventType: MessageEventTypeContentBlockStop,
			event: `{
				"type": "content_block_stop"
			}`,
			expected: &MessageStreamResponse{
				Type: "content_block_stop",
			},
		},
		{
			eventType: MessageEventTypeMessageDelta,
			event: `{
				"type": "message_delta",
				"delta": {
					"stop_reason": "something",
					"stop_sequence": "else"
				},
				"usage": {
					"output_tokens": 20
				}
			}`,
			expected: &MessageStreamResponse{
				Type: "message_delta",
				Delta: MessageStreamDelta{
					StopReason:   "something",
					StopSequence: "else",
				},
				Usage: MessageStreamUsage{
					OutputTokens: 20,
				},
			},
		},
		{
			eventType: MessageEventTypeMessageStop,
			event: `{
				"type": "message_stop"
			}`,
			expected: &MessageStreamResponse{
				Type: "message_stop",
			},
		},
		{
			eventType: MessageEventTypeError,
			event: `{
				"type": "message_error",
				"error": {
					"type": "error",
					"message": "This is an error"
				}
			}`,
			expected:  &MessageStreamResponse{},
			expErrStr: "error type: error, message: This is an error",
		},
	}

	for _, test := range events {
		response, err := ParseMessageEvent(test.eventType, test.event)
		if err != nil && test.expErrStr == "" {
			t.Errorf("unexpected error: %v", err)
		}

		if err != nil && err.Error() != test.expErrStr {
			t.Errorf("unexpected error, got: %v, want: %v", err, test.expErrStr)
		}

		if !reflect.DeepEqual(response, test.expected) {
			t.Errorf("unexpected response, got: %v, want: %v", response, test.expected)
		}
	}
}

func TestUnsupportedEventType(t *testing.T) {
	res, err := ParseMessageEvent(MessageEventType("not-a-real-type"), "")
	if err == nil {
		t.Errorf("expected error, got: %v", res)
	}

	if err.Error() != "unknown event type" {
		t.Errorf("unexpected error, got: %v", err)
	}
}
