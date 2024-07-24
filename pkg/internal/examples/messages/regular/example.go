package main

import (
	"context"
	"fmt"

	"github.com/madebywelch/anthropic-go/v3/pkg/anthropic"
	"github.com/madebywelch/anthropic-go/v3/pkg/anthropic/client/native"
)

func main() {
	ctx := context.Background()
	client, err := native.MakeClient(native.Config{
		APIKey: "your-api-key",
	})
	if err != nil {
		panic(err)
	}

	// Prepare a message request
	request := anthropic.NewMessageRequest(
		anthropic.WithMessageModel(anthropic.ClaudeV2_1),
		anthropic.WithMessageMaxTokens(20),
		anthropic.WithMessages([]anthropic.MessagePartRequest{{Role: "user", Content: []anthropic.ContentBlock{anthropic.NewTextContentBlock("Hello, world!")}}}),
	)

	// Call the Message method
	response, err := client.Message(ctx, request)
	if err != nil {
		panic(err)
	}

	fmt.Println(response.Content)
}
