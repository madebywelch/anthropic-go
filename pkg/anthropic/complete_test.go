package anthropic

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestComplete(t *testing.T) {
	// Create a test server to mock the Anthropics API
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"completion": "Test completion"}`)) // Mock response
	}))
	defer testServer.Close()

	// Create a new client with the test server's URL
	client, err := NewClient("fake-api-key")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	client.baseURL = testServer.URL // Override baseURL to point to the test server

	// Prepare a completion request
	request := NewCompletionRequest("Why is the sky blue?")

	// Call the Complete method
	response, err := client.Complete(request)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check the response
	expectedCompletion := "Test completion"
	if response.Completion != expectedCompletion {
		t.Errorf("Expected completion %q, got %q", expectedCompletion, response.Completion)
	}
}

func TestCompleteWithParameters(t *testing.T) {
	// Prepare a completion request
	request := NewCompletionRequest("Why is the sky blue?",
		WithModel(ClaudeInstantV1_1_100k),
		WithTemperature(0.5),
		WithMaxTokens(10),
		WithTopK(5),
		WithTopP(0.9),
		WithStopSequences([]string{"\n", "Why is the sky blue?"}),
	)

	if request.Prompt != "Why is the sky blue?" {
		t.Errorf("Expected prompt %q, got %q", "Why is the sky blue?", request.Prompt)
	}

	if request.Model != ClaudeInstantV1_1_100k {
		t.Errorf("Expected model %q, got %q", ClaudeInstantV1_1_100k, request.Model)
	}

	if request.Temperature != 0.5 {
		t.Errorf("Expected temperature %f, got %f", 0.5, request.Temperature)
	}

	if request.MaxTokensToSample != 10 {
		t.Errorf("Expected max tokens %d, got %d", 10, request.MaxTokensToSample)
	}

	if request.TopK != 5 {
		t.Errorf("Expected top k %d, got %d", 5, request.TopK)
	}

	if request.TopP != 0.9 {
		t.Errorf("Expected top p %f, got %f", 0.9, request.TopP)
	}

	if len(request.StopSequences) != 2 {
		t.Errorf("Expected stop sequences length %d, got %d", 2, len(request.StopSequences))
	}

	if request.StopSequences[0] != "\n" {
		t.Errorf("Expected stop sequence %q, got %q", "\n", request.StopSequences[0])
	}

	if request.StopSequences[1] != "Why is the sky blue?" {
		t.Errorf("Expected stop sequence %q, got %q", "Why is the sky blue?", request.StopSequences[1])
	}
}
