package anthropic

import (
	"fmt"
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
		WithModel[CompletionRequest](ClaudeInstantV1_1_100k),
		WithTemperature[CompletionRequest](0.5),
		WithMaxTokens[CompletionRequest](10),
		WithTopK[CompletionRequest](5),
		WithTopP[CompletionRequest](0.9),
		WithStopSequences[CompletionRequest]([]string{"\n", "Why is the sky blue?"}),
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

func TestCompleteIncompatibleModel(t *testing.T) {
	// Create client
	client, err := NewClient("fake-api-key")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Prepare a completion request
	request := NewCompletionRequest("Why is the sky blue?",
		WithModel[CompletionRequest](Claude3Opus),
	)

	// Call the Complete method expecting an error
	_, err = client.Complete(request)
	if err == nil {
		t.Fatal("Expected an incompatibility error, got none")
	}

	// Check the error message
	expErr := fmt.Sprintf("model %s is not compatible with the completion endpoint", request.Model)
	if err.Error() != expErr {
		t.Fatalf("Expected error %s, got %s", expErr, err.Error())
	}
}

func TestCompleteStreamNoStreamFlag(t *testing.T) {
	// Create client
	client, err := NewClient("fake-api-key")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Prepare a completion request
	request := NewCompletionRequest("Why is the sky blue?",
		WithModel[CompletionRequest](Claude3Opus),
	)

	// Call the Complete method expecting an error
	_, errCh := client.CompleteStream(request)
	err = <-errCh
	if err == nil {
		t.Fatal("Expected a missing stream flag error, got none")
	}

	// Check the error message
	expErr := "cannot use CompleteStream with a non-streaming request, use Complete instead"
	if err.Error() != expErr {
		t.Fatalf("Expected error %s, got %s", expErr, err.Error())
	}
}

func TestCompleteStreamIncompatibleModel(t *testing.T) {
	// Create client
	client, err := NewClient("fake-api-key")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Prepare a completion request
	request := NewCompletionRequest("Why is the sky blue?",
		WithModel[CompletionRequest](Claude3Opus),
		WithStream[CompletionRequest](true),
	)

	// Call the Complete method expecting an error
	_, errCh := client.CompleteStream(request)
	err = <-errCh
	if err == nil {
		t.Fatal("Expected an incompatibility error, got none")
	}

	// Check the error message
	expErr := fmt.Sprintf("model %s is not compatible with the completion endpoint", request.Model)
	if err.Error() != expErr {
		t.Fatalf("Expected error %s, got %s", expErr, err.Error())
	}
}
