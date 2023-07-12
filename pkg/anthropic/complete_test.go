package anthropic

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestComplete(t *testing.T) {
	testCases := []struct {
		name           string
		request        *CompletionRequest
		expectedStatus int
		expectedOutput *CompletionResponse
	}{
		{
			name: "valid completion request",
			request: &CompletionRequest{
				Prompt:            GetPrompt("Why is the sky blue?"),
				Model:             ClaudeV1,
				MaxTokensToSample: 10,
				Stream:            false,
			},
			expectedStatus: http.StatusOK,
			expectedOutput: &CompletionResponse{
				Completion: "The sky appears blue",
				StopReason: "stop_sequence",
				Stop:       "\n\nHuman:",
			},
		},
		{
			name: "valid streaming completion request",
			request: &CompletionRequest{
				Prompt:            GetPrompt("What is the meaning of life?"),
				Model:             ClaudeV1,
				MaxTokensToSample: 10,
				Stream:            true,
			},
			expectedStatus: http.StatusOK,
			expectedOutput: &CompletionResponse{
				Completion: "The meaning of life is",
				StopReason: "stop_sequence",
				Stop:       "\n\nHuman:",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/v1/complete" {
					http.Error(w, "invalid path", http.StatusBadRequest)
					return
				}
				if r.Method != http.MethodPost {
					http.Error(w, "invalid method", http.StatusMethodNotAllowed)
					return
				}
				if r.Header.Get("Content-Type") != "application/json" {
					http.Error(w, "invalid content type", http.StatusBadRequest)
					return
				}

				var req CompletionRequest
				err := json.NewDecoder(r.Body).Decode(&req)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				if !strings.Contains(req.Prompt, tc.request.Prompt) {
					http.Error(w, "invalid prompt", http.StatusBadRequest)
					return
				}

				if req.Stream {
					w.Header().Set("Content-Type", "text/event-stream")
					data, _ := json.Marshal(tc.expectedOutput)
					fmt.Fprintf(w, "data: %s\n\n", data)
				} else {
					w.WriteHeader(tc.expectedStatus)
					json.NewEncoder(w).Encode(tc.expectedOutput)
				}
			}))
			defer server.Close()

			client, err := NewClient("test-api-key")
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}
			client.baseURL = server.URL

			_, err = client.Complete(tc.request, func(resp *CompletionResponse) error {
				if resp == nil {
					t.Fatalf("response is nil")
				}
				if resp.Completion != tc.expectedOutput.Completion {
					t.Errorf("invalid completion: %s", resp.Completion)
				}
				if resp.StopReason != tc.expectedOutput.StopReason {
					t.Errorf("invalid stop reason: %s", resp.StopReason)
				}
				if resp.Stop != tc.expectedOutput.Stop {
					t.Errorf("invalid stop sequence: %s", resp.Stop)
				}
				return nil
			})
			if err != nil {
				t.Fatalf("failed to send request: %v", err)
			}
		})
	}
}
