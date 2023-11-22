package integration_tests

import (
	"os"
	"strings"
	"testing"

	"github.com/madebywelch/anthropic-go/pkg/anthropic"
	"github.com/madebywelch/anthropic-go/pkg/anthropic/utils"
)

func TestCompleteIntegration(t *testing.T) {
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

	// Prepare the prompt
	prompt, err := utils.GetPrompt("Why is the sky blue?")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Prepare a completion request
	request := anthropic.NewCompletionRequest(prompt)

	// Call the Complete method
	response, err := client.Complete(request)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Basic assertion to check if a completion is returned
	if response.Completion == "" {
		t.Errorf("Expected a completion, got an empty string")
	}

	t.Logf("Regular Completion: %s", response.Completion)
}

func TestCompleteStreamIntegration(t *testing.T) {
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

	// Prepare the prompt
	prompt, err := utils.GetPrompt("Why is the sky blue?")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Prepare a completion request
	request := anthropic.NewCompletionRequest(prompt, anthropic.WithStreaming(true), anthropic.WithMaxTokens(10), anthropic.WithModel(anthropic.ClaudeV2_1))

	// Call the Complete method (should return an error since streaming is enabled)
	_, err = client.Complete(request)
	if err == nil {
		t.Fatalf("Expected error: %v", err)
	}

	// Call the CompleteStream method
	res, errs := client.CompleteStream(request)

	MAX_ITERATIONS := 10
	builder := strings.Builder{}

main:
	for {
		select {
		case err := <-errs:
			t.Fatalf("Unexpected error: %v", err)
		case event := <-res:
			t.Logf("Completion: %s", event.Completion)
			builder.WriteString(event.Completion)
			MAX_ITERATIONS--
			if MAX_ITERATIONS == 0 {
				break main
			}
		}
	}

	t.Logf("Stream Completion: %s", builder.String())
}
