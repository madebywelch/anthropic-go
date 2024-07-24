package integration_tests

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/madebywelch/anthropic-go/v3/pkg/anthropic"
	"github.com/madebywelch/anthropic-go/v3/pkg/anthropic/client/native"
	"github.com/madebywelch/anthropic-go/v3/pkg/anthropic/utils"
)

func TestCompleteIntegration(t *testing.T) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		t.Skip("ANTHROPIC_API_KEY environment variable is not set, skipping integration test")
	}

	ctx := context.Background()

	anthropicClient, err := native.MakeClient(native.Config{
		APIKey: apiKey,
	})
	if err != nil {
		t.Fatalf("Unexpected error creating client: %v", err)
	}

	testCases := []struct {
		name         string
		model        anthropic.Model
		prompt       string
		options      []anthropic.CompletionRequestOption
		stopSequence string
	}{
		{
			name:   "Claude-2 Default",
			model:  anthropic.ClaudeV2,
			prompt: "Explain quantum computing in simple terms.",
		},
		{
			name:   "Claude-2.1 with Custom Options",
			model:  anthropic.ClaudeV2_1,
			prompt: "Write a haiku about artificial intelligence.",
			options: []anthropic.CompletionRequestOption{
				anthropic.WithCompletionMaxTokens(100),
				anthropic.WithCompletionTemperature(0.7),
				anthropic.WithCompletionTopP(0.9),
			},
		},
		{
			name:         "Claude-Instant-1 with Stop Sequence",
			model:        anthropic.ClaudeInstantV1,
			prompt:       "List the first 5 prime numbers:",
			stopSequence: ".",
			options: []anthropic.CompletionRequestOption{
				anthropic.WithCompletionStopSequences([]string{"."}),
			},
		},
		{
			name:   "Claude-2 with Low Temperature",
			model:  anthropic.ClaudeV2,
			prompt: "What is the capital of France?",
			options: []anthropic.CompletionRequestOption{
				anthropic.WithCompletionTemperature(0.1),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Prepare the prompt
			prompt, err := utils.GetPrompt(tc.prompt)
			if err != nil {
				t.Fatalf("Unexpected error preparing prompt: %v", err)
			}

			options := append([]anthropic.CompletionRequestOption{
				anthropic.WithCompletionModel(tc.model),
			}, tc.options...)
			request := anthropic.NewCompletionRequest(prompt, options...)

			response, err := anthropicClient.Complete(ctx, request)
			if err != nil {
				t.Fatalf("Unexpected error in Complete: %v", err)
			}

			if response.Completion == "" {
				t.Errorf("Expected a completion, got an empty string")
			}

			t.Logf("Model: %s, Prompt: %q", tc.model, tc.prompt)
			t.Logf("Completion: %s", response.Completion)

			switch tc.model {
			case anthropic.ClaudeInstantV1:
				if len(response.Completion) > 100 {
					t.Logf("Warning: Claude Instant-1 response longer than expected: %d characters", len(response.Completion))
				}
			case anthropic.ClaudeV2, anthropic.ClaudeV2_1:
				if len(response.Completion) < 50 {
					t.Logf("Warning: Claude-2/2.1 response shorter than expected: %d characters", len(response.Completion))
				}
			}

			if tc.stopSequence != "" {
				if strings.Contains(response.Completion, tc.stopSequence) {
					t.Errorf("Completion contains stop sequence %q", tc.stopSequence)
				}
			}
		})
	}
}
