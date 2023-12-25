package integration_tests

import (
	"os"
	"testing"

	"github.com/madebywelch/anthropic-go/v2/pkg/anthropic"
)

func TestMessageIntegration(t *testing.T) {
	// Get the API key from the environment
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		t.Skip("ANTHROPIC_API_KEY environment variable is not set, skipping integration test")
	}

	// Create a new client
	client, err := anthropic.NewClient(apiKey)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Prepare a message request
	request := &anthropic.MessageRequest{
		Model:             anthropic.ClaudeV2_1,
		MaxTokensToSample: 10,
		Messages:          []anthropic.MessagePartRequest{{Role: "user", Content: "Hello, Anthropics!"}},
	}

	// Call the Message method
	response, err := client.Message(request)
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
	client, err := anthropic.NewClient(apiKey)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Prepare a message request
	request := &anthropic.MessageRequest{
		Model:    anthropic.ClaudeV2_1,
		Messages: []anthropic.MessagePartRequest{{Role: "user", Content: "Hello, Anthropics!"}},
	}

	// Call the Message method expecting an error
	_, err = client.Message(request)
	// We're expecting an error here because we didn't set the required field MaxTokensToSample
	if err == nil {
		t.Fatal("Expected an error, got none")
	}
}

// - TODO: TestMessageWithParametersIntegration: to test sending a message with various parameters
// - TODO: TestMessageStreamIntegration: to ensure the function correctly handles streaming requests
