package integration_tests

import (
	"context"
	"os"
	"testing"

	"github.com/madebywelch/anthropic-go/v3/pkg/anthropic"
	"github.com/madebywelch/anthropic-go/v3/pkg/anthropic/client/native"
)

func TestMessageWithToolsIntegration(t *testing.T) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		t.Skip("ANTHROPIC_API_KEY environment variable is not set, skipping integration test")
	}

	ctx := context.Background()

	anthropicClient, err := native.MakeClient(native.Config{
		APIKey: apiKey,
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	request := &anthropic.MessageRequest{
		Model:             anthropic.Claude3Opus,
		MaxTokensToSample: 1024,
		Tools: []anthropic.Tool{
			{
				Name:        "get_weather",
				Description: "Get the weather",
				InputSchema: anthropic.InputSchema{
					Type: "object",
					Properties: map[string]anthropic.Property{
						"city": {Type: "string", Description: "city to get the weather for"},
						"unit": {Type: "string", Enum: []string{"celsius", "fahrenheit"}, Description: "temperature unit to return"}},
					Required: []string{"city"},
				},
			},
		},
		Messages: []anthropic.MessagePartRequest{
			{
				Role: "user",
				Content: []anthropic.ContentBlock{
					anthropic.NewTextContentBlock("what's the weather in Charleston?"),
				},
			},
		},
	}

	response, err := anthropicClient.Message(ctx, request)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if response == nil || len(response.Content) == 0 {
		t.Errorf("Expected a message response, got none or empty content")
	}

	if response.StopReason != "tool_use" {
		t.Errorf("Expected stop reason 'tool_use', got %s", response.StopReason)
	}
}

func TestMessageWithForcedToolIntegration(t *testing.T) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		t.Skip("ANTHROPIC_API_KEY environment variable is not set, skipping integration test")
	}

	anthropicClient, err := native.MakeClient(native.Config{
		APIKey: apiKey,
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	request := &anthropic.MessageRequest{
		Model:             anthropic.Claude3Opus,
		MaxTokensToSample: 1024,
		ToolChoice: &anthropic.ToolChoice{
			Type: "tool",
			Name: "get_weather",
		},
		Tools: []anthropic.Tool{
			{
				Name:        "get_weather",
				Description: "Get the weather",
				InputSchema: anthropic.InputSchema{
					Type: "object",
					Properties: map[string]anthropic.Property{
						"city": {Type: "string", Description: "city to get the weather for"},
						"unit": {Type: "string", Enum: []string{"celsius", "fahrenheit"}, Description: "temperature unit to return"}},
					Required: []string{"city"},
				},
			},
		},
		Messages: []anthropic.MessagePartRequest{
			{
				Role: "user",
				Content: []anthropic.ContentBlock{
					anthropic.NewTextContentBlock("what's the weather in Charleston?"),
				},
			},
		},
	}

	response, err := anthropicClient.Message(context.Background(), request)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if response == nil || len(response.Content) == 0 {
		t.Errorf("Expected a message response, got none or empty content")
	}

	if response.StopReason != "tool_use" {
		t.Errorf("Expected stop reason 'tool_use', got %s", response.StopReason)
	}
}

func TestMessageWithImageIntegration(t *testing.T) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		t.Skip("ANTHROPIC_API_KEY environment variable is not set, skipping integration test")
	}

	anthropicClient, err := native.MakeClient(native.Config{
		APIKey: apiKey,
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	request := &anthropic.MessageRequest{
		Model:             anthropic.Claude3Opus,
		MaxTokensToSample: 50,
		Messages: []anthropic.MessagePartRequest{
			{
				Role: "user",
				Content: []anthropic.ContentBlock{
					anthropic.NewImageContentBlock(anthropic.MediaTypePNG, "iVBORw0KGgoAAAANSUhEUgAAAAoAAAAKCAYAAACNMs+9AAAAFUlEQVR42mP8z8BQz0AEYBxVSF+FABJADveWkH6oAAAAAElFTkSuQmCC"),
					anthropic.NewTextContentBlock("What is this image?"),
				},
			},
		},
	}

	response, err := anthropicClient.Message(context.Background(), request)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if response == nil || len(response.Content) == 0 {
		t.Errorf("Expected a message response, got none or empty content")
	}
}

func TestMessageIntegration(t *testing.T) {
	// Get the API key from the environment
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		t.Skip("ANTHROPIC_API_KEY environment variable is not set, skipping integration test")
	}

	// Create a new client
	anthropicClient, err := native.MakeClient(native.Config{
		APIKey: apiKey,
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Prepare a message request
	request := &anthropic.MessageRequest{
		Model:             anthropic.ClaudeV2_1,
		MaxTokensToSample: 10,
		Messages:          []anthropic.MessagePartRequest{{Role: "user", Content: []anthropic.ContentBlock{anthropic.NewTextContentBlock("Hello, Anthropics!")}}},
	}

	// Call the Message method
	response, err := anthropicClient.Message(context.Background(), request)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Basic assertion to check if a message response is returned
	if response == nil || len(response.Content) == 0 {
		t.Errorf("Expected a message response, got none or empty content")
	}

	// Ensure the response contains populated ID
	if response.ID == "" {
		t.Errorf("Expected a message response with a non-empty ID, got none")
	}
}

func TestMessageErrorHandlingIntegration(t *testing.T) {
	// Get the API key from the environment
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		t.Skip("ANTHROPIC_API_KEY environment variable is not set, skipping integration test")
	}

	// Create a new client
	anthropicClient, err := native.MakeClient(native.Config{
		APIKey: apiKey,
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Prepare a message request
	request := &anthropic.MessageRequest{
		Model:    anthropic.ClaudeV2_1,
		Messages: []anthropic.MessagePartRequest{{Role: "user", Content: []anthropic.ContentBlock{anthropic.NewTextContentBlock("Hello, Anthropics!")}}},
	}

	// Call the Message method expecting an error
	_, err = anthropicClient.Message(context.Background(), request)
	// We're expecting an error here because we didn't set the required field MaxTokensToSample
	if err == nil {
		t.Fatal("Expected an error, got none")
	}
}
