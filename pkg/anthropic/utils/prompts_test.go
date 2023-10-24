package utils

import (
	"testing"
)

func TestGetPrompt(t *testing.T) {
	got, err := GetPrompt("Hello, Claude")
	if err != nil {
		t.Errorf("GetPrompt(\"Hello, Claude\") returned an error: %v", err)
	}
	want := "\n\nHuman: Hello, Claude\n\nAssistant:"
	if got != want {
		t.Errorf("GetPrompt(\"Hello, Claude\") = %s; want %s", got, want)
	}
}

func TestGetChatPrompt(t *testing.T) {
	messages := []Message{
		{Sender: "Human", Content: "Hello, Claude"},
		{Sender: "Assistant", Content: "Hello, Human"},
	}

	got, err := GetChatPrompt(messages)
	if err != nil {
		t.Errorf("GetChatPrompt(...) returned an error: %v", err)
	}
	want := "\n\nHuman: Hello, Claude\n\nAssistant: Hello, Human\n\nAssistant:"
	if got != want {
		t.Errorf("GetChatPrompt(...) = %s; want %s", got, want)
	}
}
