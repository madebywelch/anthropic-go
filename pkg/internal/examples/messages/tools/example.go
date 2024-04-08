package main

import (
	"github.com/madebywelch/anthropic-go/v2/pkg/anthropic"
)

func main() {
	client, err := anthropic.NewClient("your-api-key")
	if err != nil {
		panic(err)
	}

	// Prepare a message request
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

	// Call the Message method
	response, err := client.Message(request)
	if err != nil {
		panic(err)
	}

	if response.StopReason == "tool_use" {
		// Do something with the tool response
	}
}
