package integration_tests

/*
func TestMessageWithToolsIntegration(t *testing.T) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		t.Skip("ANTHROPIC_API_KEY environment variable is not set, skipping integration test")
	}

	client, err := anthropic.NewClient(apiKey)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	request := &anthropic.MessageRequest{
		Model:             anthropic.Claude3Opus,
		MaxTokensToSample: 1024,
		Tools: []anthropic.Tool{
			{
				Name:        "get_weather",
				Description: "Get the weather",
				InputSchema: anthropic.InputSchema{
					Type: "object",
					Properties: map[string]anthropic.Property{
						"city": {Type: "string", Description: "city to get the weather for"},
						"unit": {Type: "string", Enum: []string{"celsius", "fahrenheit"}, Description: "temperature unit to return"}},
					Required: []string{"city"},
				},
			},
		},
		Messages: []anthropic.MessagePartRequest{
			{
				Role: "user",
				Content: []anthropic.ContentBlock{
					anthropic.NewTextContentBlock("what's the weather in Charleston?"),
				},
			},
		},
	}

	response, err := client.Message(request)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if response == nil || len(response.Content) == 0 {
		t.Errorf("Expected a message response, got none or empty content")
	}

	if response.StopReason != "tool_use" {
		t.Errorf("Expected stop reason 'tool_use', got %s", response.StopReason)
	}
}
*/
