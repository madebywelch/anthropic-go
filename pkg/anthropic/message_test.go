package anthropic

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMessage(t *testing.T) {
	// Mock server for successful message response
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := MessageResponse{
			ID:    "12345",
			Type:  "testType",
			Model: "testModel",
			Role:  "user",
			Content: []MessagePartResponse{{
				Type: "text",
				Text: "Test message",
			}},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer testServer.Close()

	// Create a new client with the test server's URL
	client, err := NewClient("fake-api-key")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	client.baseURL = testServer.URL // Override baseURL to point to the test server

	// Prepare a message request
	request := &MessageRequest{
		Model:    ClaudeV2_1,
		Messages: []MessagePartRequest{{Role: "user", Content: "Hello"}},
	}

	// Call the Message method
	response, err := client.Message(request)
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
	client, err := NewClient("fake-api-key")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	client.baseURL = testServer.URL // Override baseURL to point to the test server

	// Prepare a message request
	request := &MessageRequest{
		Model:    ClaudeV2_1,
		Messages: []MessagePartRequest{{Role: "user", Content: "Hello"}},
	}

	// Call the Message method expecting an error
	_, err = client.Message(request)
	if err == nil {
		t.Fatal("Expected an error, got none")
	}
}

func TestMessageStreamNotSupported(t *testing.T) {
	// Create client
	client, err := NewClient("fake-api-key")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Prepare a message request with streaming set to true
	request := &MessageRequest{
		Model:    ClaudeV2_1,
		Messages: []MessagePartRequest{{Role: "user", Content: "Hello"}},
		Stream:   true,
	}

	// Call the Message method expecting an error
	_, err = client.Message(request)
	if err == nil {
		t.Fatal("Expected an error for streaming not supported, got none")
	}
}
