package integration_tests

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/madebywelch/anthropic-go/v3/pkg/anthropic"
	"github.com/madebywelch/anthropic-go/v3/pkg/anthropic/client/native"
)

func TestToolIntegration(t *testing.T) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		t.Skip("ANTHROPIC_API_KEY environment variable is not set, skipping integration test")
	}

	client, err := native.MakeClient(native.Config{APIKey: apiKey})
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	t.Run("Weather Tool Usage", func(t *testing.T) {
		request := createWeatherRequest()
		response, err := client.Message(context.Background(), request)
		if err != nil {
			t.Fatalf("Error sending message: %v", err)
		}

		toolUseID, assistantMessageBlock := processResponse(t, response)
		request.AddAssistantMessage(assistantMessageBlock...)

		request.Messages = append(request.Messages, createToolResultMessage(toolUseID))

		if len(request.Messages) != len(response.Content)+1 {
			t.Errorf("Expected %d messages, got %d", len(response.Content)+1, len(request.Messages))
		}

		finalResponse, err := client.Message(context.Background(), request)
		if err != nil {
			t.Fatalf("Error sending final message: %v", err)
		}

		if !strings.Contains(finalResponse.Content[0].Text, "52") {
			t.Errorf("Response should contain the temperature '52', got: %s", finalResponse.Content[0].Text)
		}
		t.Logf("Tool Response: %s", finalResponse.Content[0].Text)
	})
}

func createWeatherRequest() *anthropic.MessageRequest {
	return &anthropic.MessageRequest{
		Model:             anthropic.Claude35Sonnet,
		MaxTokensToSample: 512,
		Tools: []anthropic.Tool{
			{
				Name:        "get_weather",
				Description: "Get the current weather in a given location",
				InputSchema: anthropic.GenerateInputSchema(&WeatherRequest{}),
			},
		},
		Messages: []anthropic.MessagePartRequest{
			{
				Role: "user",
				Content: []anthropic.ContentBlock{
					anthropic.NewTextContentBlock("What is the weather in fahrenheit like in Charleston, SC?"),
				},
			},
		},
	}
}

func processResponse(t *testing.T, response *anthropic.MessageResponse) (string, []anthropic.ContentBlock) {
	t.Helper()
	var toolUseID string
	assistantMessageBlock := []anthropic.ContentBlock{anthropic.NewTextContentBlock(response.Content[0].Text)}

	for _, part := range response.Content {
		if part.Type == "tool_use" {
			toolUseID = part.ID
			assistantMessageBlock = append(assistantMessageBlock, anthropic.NewToolUseContentBlock(part.ID, "get_weather", part.Input))
		}
	}

	if toolUseID == "" {
		t.Fatal("Tool Use ID is empty")
	}
	return toolUseID, assistantMessageBlock
}

func createToolResultMessage(toolUseID string) anthropic.MessagePartRequest {
	return anthropic.MessagePartRequest{
		Role: "user",
		Content: []anthropic.ContentBlock{
			anthropic.NewToolResultContentBlock(toolUseID, "the temperature is 52f", false),
		},
	}
}
