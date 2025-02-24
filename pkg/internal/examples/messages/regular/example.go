package main

import (
	"context"
	"fmt"

	"github.com/pigeonlaser/anthropic-go/v3/pkg/anthropic"
	"github.com/pigeonlaser/anthropic-go/v3/pkg/anthropic/client/native"
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
		[]anthropic.MessagePartRequest{{Role: "user", Content: []anthropic.ContentBlock{anthropic.NewTextContentBlock("Hello, world!")}}},
		anthropic.WithModel[anthropic.MessageRequest](anthropic.ClaudeV2_1),
		anthropic.WithMaxTokens[anthropic.MessageRequest](20),
	)

	// Call the Message method
	response, err := client.Message(ctx, request)
	if err != nil {
		panic(err)
	}

	fmt.Println(response.Content)
}
