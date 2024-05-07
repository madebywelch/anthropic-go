package native

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/madebywelch/anthropic-go/v3/pkg/anthropic"
)

func TestMessageStreamNoStreamFlag(t *testing.T) {
	// Create client
	client, err := MakeClient(Config{
		APIKey: "fake-api-key",
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Prepare a message request without streaming
	request := &anthropic.MessageRequest{
		Model: anthropic.Claude3Opus,
		Messages: []anthropic.MessagePartRequest{{
			Role:    "user",
			Content: []anthropic.ContentBlock{anthropic.NewTextContentBlock("Hello")},
		}},
	}

	// Call the MessageStream method expecting an error
	_, errCh := client.MessageStream(context.Background(), request)

	err = <-errCh
	if err == nil {
		t.Fatal("Expected an error for streaming without a stream request")
	}

	expErr := "cannot use MessageStream with streaming disabled, use Message instead"

	if err.Error() != expErr {
		t.Fatalf(
			"Expected error %s, got %s",
			expErr,
			err.Error(),
		)
	}
}

func TestMessageStreamIncompatibleModel(t *testing.T) {
	// Create client
	client, err := MakeClient(Config{
		APIKey: "fake-api-key",
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Prepare a message request without streaming
	request := &anthropic.MessageRequest{
		Model: anthropic.ClaudeV2,
		Messages: []anthropic.MessagePartRequest{{
			Role:    "user",
			Content: []anthropic.ContentBlock{anthropic.NewTextContentBlock("Hello")},
		}},
		Stream: true,
	}

	// Call the MessageStream method expecting an error
	_, errCh := client.MessageStream(context.Background(), request)

	err = <-errCh
	if err == nil {
		t.Fatal("Expected an error for streaming not supported, got none")
	}

	expErr := fmt.Sprintf("model %s is not compatible with the messagestream endpoint", request.Model)

	if err.Error() != expErr {
		t.Fatalf(
			"Expected error %s, got %s",
			expErr,
			err.Error(),
		)
	}
}

func TestMessageStreamSuccess(t *testing.T) {
	// Create a test server to mock the Anthropics API
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")

		output := []byte{}
		output = append(output, []byte("event: message_start\n")...)
		output = append(output, []byte("data: {\"type\": \"message_start\", \"message\": {\"id\": \"msg_1nZdL29xx5MUA1yADyHTEsnR8uuvGzszyY\", \"type\": \"message\", \"role\": \"assistant\", \"content\": [], \"model\": \"claude-3-opus-20240229\", \"stop_reason\": null, \"stop_sequence\": null, \"usage\": {\"input_tokens\": 25, \"output_tokens\": 1}}}\n\n")...)
		output = append(output, []byte("event: content_block_start\n")...)
		output = append(output, []byte("data: {\"type\": \"content_block_start\", \"index\":0, \"content_block\": {\"type\": \"text\", \"text\": \"\"}}\n\n")...)
		output = append(output, []byte("event: ping\n")...)
		output = append(output, []byte("data: {\"type\": \"ping\"}\n\n")...)
		output = append(output, []byte("event: not_really_anything_just_testing_unknown_event\n")...)
		output = append(output, []byte("data: {\"type\": \"not_really_anything_just_testing_unknown_event\"}\n\n")...)
		output = append(output, []byte("event: content_block_delta\n")...)
		output = append(output, []byte("data: {\"type\": \"content_block_delta\", \"index\": 0, \"delta\": {\"type\": \"text_delta\", \"text\": \"hello there, \"}}\n\n")...)
		output = append(output, []byte("event: content_block_delta\n")...)
		output = append(output, []byte("data: {\"type\": \"content_block_delta\", \"index\": 0, \"delta\": {\"type\": \"text_delta\", \"text\": \"general kenobi\"}}\n\n")...)
		output = append(output, []byte("event: content_block_stop\n")...)
		output = append(output, []byte("data: {\"type\": \"content_block_stop\", \"index\": 0}\n\n")...)
		output = append(output, []byte("event: message_delta\n")...)
		output = append(output, []byte("data: {\"type\": \"message_delta\", \"delta\": {\"stop_reason\": \"end_turn\", \"stop_sequence\":null}, \"usage\":{\"output_tokens\": 15}}\n\n")...)
		output = append(output, []byte("event: message_stop\n")...)
		output = append(output, []byte("data: {\"type\": \"message_stop\"}\n\n")...)
		w.Write(output)
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

	// Prepare a message request
	request := &anthropic.MessageRequest{
		Model: anthropic.Claude3Opus,
		Messages: []anthropic.MessagePartRequest{{
			Role:    "user",
			Content: []anthropic.ContentBlock{anthropic.NewTextContentBlock("Hello")},
		}},
		Stream: true,
	}

	rCh, errCh := client.MessageStream(context.Background(), request)
	chunk := &anthropic.MessageStreamResponse{}
	final := strings.Builder{}
	inputTokens := 0
	outputTokens := 0

	for {
		select {
		case chunk = <-rCh:
			final.WriteString(chunk.Delta.Text)
			if chunk.Type == "message_start" {
				inputTokens = chunk.Usage.InputTokens
			} else if chunk.Type == "message_delta" {
				outputTokens = chunk.Usage.OutputTokens
			}
		case err := <-errCh:
			t.Fatalf("Unexpected error: %s", err.Error())
		}
		if chunk.Type == "message_stop" {
			break
		}
	}

	// Check the response
	expectedResult := "hello there, general kenobi"
	if final.String() != expectedResult {
		t.Fatalf("Expected result %s, got %s", expectedResult, final.String())
	}

	// Check the usage
	if inputTokens != 25 {
		t.Fatalf("Expected input tokens %d, got %d", 25, inputTokens)
	}

	if outputTokens != 15 {
		t.Fatalf("Expected output tokens %d, got %d", 15, outputTokens)
	}
}

func TestMessageStreamErrorInStream(t *testing.T) {
	// Create a test server to mock the Anthropics API
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")

		output := []byte{}
		output = append(output, []byte("event: message_start\n")...)
		output = append(output, []byte("data: {\"type\": \"message_start\", \"message\": {\"id\": \"msg_1nZdL29xx5MUA1yADyHTEsnR8uuvGzszyY\", \"type\": \"message\", \"role\": \"assistant\", \"content\": [], \"model\": \"claude-3-opus-20240229\", \"stop_reason\": null, \"stop_sequence\": null, \"usage\": {\"input_tokens\": 25, \"output_tokens\": 1}}}\n\n")...)
		output = append(output, []byte("event: content_block_start\n")...)
		output = append(output, []byte("data: {\"type\": \"content_block_start\", \"index\":0, \"content_block\": {\"type\": \"text\", \"text\": \"\"}}\n\n")...)
		output = append(output, []byte("event: error\n")...)
		output = append(output, []byte("data: {\"type\": \"error\", \"error\": {\"type\": \"overload_error\", \"message\": \"Overloaded\"}}\n\n")...)
		w.Write(output)
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

	// Prepare a message request
	request := &anthropic.MessageRequest{
		Model: anthropic.Claude3Opus,
		Messages: []anthropic.MessagePartRequest{{
			Role:    "user",
			Content: []anthropic.ContentBlock{anthropic.NewTextContentBlock("Hello")},
		}},
		Stream: true,
	}

	// Call the Complete method
	rCh, errCh := client.MessageStream(context.Background(), request)
	var chunk *anthropic.MessageStreamResponse
	final := strings.Builder{}
	done := false

	for {
		select {
		case chunk = <-rCh:
			final.WriteString(chunk.Delta.Text)
		case err = <-errCh:
			done = true
			break
		}
		if chunk.Type == "message_stop" || done {
			break
		}
	}

	// Check the response is empty
	expectedResult := ""
	if final.String() != expectedResult {
		t.Fatalf("Expected result %s, got %s", expectedResult, final.String())
	}

	// Check the error
	expectedError := "error processing message stream: error type: overload_error, message: Overloaded"
	if expectedError != err.Error() {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}
