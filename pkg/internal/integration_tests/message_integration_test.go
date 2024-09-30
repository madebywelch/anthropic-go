package integration_tests

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/madebywelch/anthropic-go/v4/pkg/anthropic"
	"github.com/madebywelch/anthropic-go/v4/pkg/anthropic/client/native"
)

type WeatherRequest struct {
	City string `json:"city" jsonschema:"required,description=city to get the weather for"`
	Unit string `json:"unit" jsonschema:"enum=celsius,enum=fahrenheit,description=temperature unit to return"`
}

func TestMessageWithToolsIntegration(t *testing.T) {
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
			Type: "auto",
		},
		Tools: []anthropic.Tool{
			{
				Name:        "get_weather",
				Description: "Get the weather",
				InputSchema: anthropic.GenerateInputSchema(&WeatherRequest{}),
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
				InputSchema: anthropic.GenerateInputSchema(&WeatherRequest{}),
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

func TestMessageIntegration(t *testing.T) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		t.Skip("ANTHROPIC_API_KEY environment variable is not set, skipping integration test")
	}

	anthropicClient, err := native.MakeClient(native.Config{
		APIKey: apiKey,
	})
	if err != nil {
		t.Fatalf("Unexpected error creating client: %v", err)
	}

	testCases := []struct {
		name    string
		request *anthropic.MessageRequest
		check   func(*testing.T, *anthropic.MessageResponse, error)
	}{
		{
			name: "Basic Message",
			request: &anthropic.MessageRequest{
				Model:             anthropic.Claude3Opus,
				MaxTokensToSample: 100,
				Messages: []anthropic.MessagePartRequest{
					{Role: "user", Content: []anthropic.ContentBlock{anthropic.NewTextContentBlock("What is the capital of France?")}},
				},
			},
			check: func(t *testing.T, response *anthropic.MessageResponse, err error) {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				if !strings.Contains(response.Content[0].Text, "Paris") {
					t.Errorf("Expected response to mention Paris, got: %s", response.Content[0].Text)
				}
			},
		},
		{
			name: "Multi-turn Conversation",
			request: &anthropic.MessageRequest{
				Model:             anthropic.Claude3Opus,
				MaxTokensToSample: 100,
				Messages: []anthropic.MessagePartRequest{
					{Role: "user", Content: []anthropic.ContentBlock{anthropic.NewTextContentBlock("What's the largest planet in our solar system?")}},
					{Role: "assistant", Content: []anthropic.ContentBlock{anthropic.NewTextContentBlock("The largest planet in our solar system is Jupiter.")}},
					{Role: "user", Content: []anthropic.ContentBlock{anthropic.NewTextContentBlock("What's the second largest?")}},
				},
			},
			check: func(t *testing.T, response *anthropic.MessageResponse, err error) {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				if !strings.Contains(response.Content[0].Text, "Saturn") {
					t.Errorf("Expected response to mention Saturn, got: %s", response.Content[0].Text)
				}
			},
		},
		{
			name: "System Prompt",
			request: &anthropic.MessageRequest{
				Model:             anthropic.Claude3Opus,
				MaxTokensToSample: 100,
				SystemPrompt:      "You are a helpful assistant that always responds in rhyme.",
				Messages: []anthropic.MessagePartRequest{
					{Role: "user", Content: []anthropic.ContentBlock{anthropic.NewTextContentBlock("Tell me about the weather.")}},
				},
			},
			check: func(t *testing.T, response *anthropic.MessageResponse, err error) {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				if !strings.Contains(response.Content[0].Text, "\n") {
					t.Errorf("Expected response to be in rhyme (contain newlines), got: %s", response.Content[0].Text)
				}
			},
		},
		{
			name: "Temperature and Top P",
			request: &anthropic.MessageRequest{
				Model:             anthropic.Claude3Opus,
				MaxTokensToSample: 100,
				Temperature:       0.9,
				TopP:              0.95,
				Messages: []anthropic.MessagePartRequest{
					{Role: "user", Content: []anthropic.ContentBlock{anthropic.NewTextContentBlock("Generate a random word.")}},
				},
			},
			check: func(t *testing.T, response *anthropic.MessageResponse, err error) {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				// We can't check for randomness easily, but we can ensure a response was received
				if response.Content[0].Text == "" {
					t.Errorf("Expected a non-empty response")
				}
			},
		},
		{
			name: "Stop Sequences",
			request: &anthropic.MessageRequest{
				Model:             anthropic.Claude3Opus,
				MaxTokensToSample: 100,
				StopSequences:     []string{".", "!"},
				Messages: []anthropic.MessagePartRequest{
					{Role: "user", Content: []anthropic.ContentBlock{anthropic.NewTextContentBlock("Write a short sentence")}},
				},
			},
			check: func(t *testing.T, response *anthropic.MessageResponse, err error) {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				if strings.ContainsAny(response.Content[0].Text, ".!") {
					t.Errorf("Expected response without '.' or '!', got: %s", response.Content[0].Text)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response, err := anthropicClient.Message(context.Background(), tc.request)
			tc.check(t, response, err)
		})
	}
}
