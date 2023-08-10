package anthropic

import (
	"fmt"
	"strings"
)

// Message represents a single message in a chat conversation.
type Message struct {
	Sender  string // The sender's name (e.g., "Human" or a username)
	Content string // The content of the message
}

// GetPrompt returns a prompt string that can be used to complete a user question.
func GetPrompt(userQuestion string) string {
	return fmt.Sprintf("Human: %s\n\nAssistant:", strings.TrimSpace(userQuestion))
}

// GetChatPrompt constructs a prompt string from a series of messages in a chat conversation.
func GetChatPrompt(chat []Message) string {
	var builder strings.Builder
	for _, message := range chat {
		builder.WriteString(fmt.Sprintf("%s: %s\n\n", message.Sender, message.Content))
	}
	builder.WriteString("Assistant:")
	return builder.String()
}
