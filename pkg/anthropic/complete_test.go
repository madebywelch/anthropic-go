package anthropic

import (
	"encoding/json"
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
			},
			expectedStatus: http.StatusOK,
			expectedOutput: &CompletionResponse{
				Completion: "The sky appears blue",
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
				if !strings.Contains(req.Prompt, "Why is the sky blue?") {
					http.Error(w, "invalid prompt", http.StatusBadRequest)
					return
				}

				w.WriteHeader(tc.expectedStatus)
				json.NewEncoder(w).Encode(tc.expectedOutput)
			}))
			defer server.Close()

			client, err := NewClient("test-api-key", WithMaxRetries(1), WithRetryDelay(10))
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}
			client.baseURL = server.URL

			resp, err := client.Complete(tc.request)
			if err != nil {
				t.Fatalf("failed to send request: %v", err)
			}

			if resp == nil {
				t.Fatalf("response is nil")
			}
			if resp.Completion != "The sky appears blue" {
				t.Errorf("invalid completion: %s", resp.Completion)
			}
			if resp.StopReason != "stop_sequence" {
				t.Errorf("invalid stop reason: %s", resp.StopReason)
			}
			if resp.Stop != "\n\nHuman:" {
				t.Errorf("invalid stop sequence: %s", resp.Stop)
			}
		})
	}
}
