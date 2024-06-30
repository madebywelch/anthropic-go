package native

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/madebywelch/anthropic-go/v3/pkg/anthropic"
)

func TestMessage(t *testing.T) {
	// Mock server for successful message response
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := &anthropic.MessageResponse{
			ID:    "12345",
			Type:  "testType",
			Model: "testModel",
			Role:  "user",
			Content: []anthropic.MessagePartResponse{{
				Type: "text",
				Text: "Test message",
			}},
			Usage: anthropic.MessageUsage{
				InputTokens:  10,
				OutputTokens: 5,
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer testServer.Close()

	ctx := context.Background()

	// Create a new client with the test server's URL
	client, err := MakeClient(Config{
		APIKey:  "fake-api-key",
		BaseURL: testServer.URL,
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Prepare a message request
	request := &anthropic.MessageRequest{
		Model: anthropic.Claude3Opus,
		Messages: []anthropic.MessagePartRequest{{
			Role:    "user",
			Content: []anthropic.ContentBlock{anthropic.NewTextContentBlock("Hello")},
		}},
	}

	// Call the Message method
	response, err := client.Message(ctx, request)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check the response
	expectedContent := "Test message"
	if len(response.Content) == 0 || response.Content[0].Text != expectedContent {
		t.Errorf("Expected message %q, got %q", expectedContent, response.Content[0].Text)
	}
}

func TestMessageErrorHandling(t *testing.T) {
	// Mock server for error response
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}))
	defer testServer.Close()

	// Create a new client with the test server's URL
	client, err := MakeClient(Config{
		APIKey:  "fake-api-key",
		BaseURL: testServer.URL,
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	client.baseURL = testServer.URL // Override baseURL to point to the test server

	// Prepare a message request
	request := &anthropic.MessageRequest{
		Model: anthropic.Claude3Opus,
		Messages: []anthropic.MessagePartRequest{{
			Role:    "user",
			Content: []anthropic.ContentBlock{anthropic.NewTextContentBlock("Hello")},
		}},
	}

	// Call the Message method expecting an error
	_, err = client.Message(context.Background(), request)
	if err == nil {
		t.Fatal("Expected an error, got none")
	}
}

func TestMessageIncompatibleModel(t *testing.T) {
	// Create client
	client, err := MakeClient(Config{APIKey: "fake-api-key"})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Prepare a message request with streaming set to true
	request := &anthropic.MessageRequest{
		Model: anthropic.ClaudeV2,
		Messages: []anthropic.MessagePartRequest{{
			Role:    "user",
			Content: []anthropic.ContentBlock{anthropic.NewTextContentBlock("Hello")},
		}},
	}

	// Call the MessageStream method expecting an error
	_, err = client.Message(context.Background(), request)

	if err == nil {
		t.Fatal("Expected an error for streaming not supported, got none")
	}

	expErr := fmt.Sprintf("model %s is not compatible with the message endpoint", request.Model)

	if err.Error() != expErr {
		t.Fatalf(
			"Expected error %s, got %s",
			expErr,
			err.Error(),
		)
	}
}
