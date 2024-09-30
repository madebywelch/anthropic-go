package main

import (
	"context"
	"fmt"

	"github.com/madebywelch/anthropic-go/v4/pkg/anthropic"
	"github.com/madebywelch/anthropic-go/v4/pkg/anthropic/client/native"
)

func main() {
	ctx := context.Background()
	client, err := native.MakeClient(native.Config{
		APIKey: "your-api-key",
	})
	if err != nil {
		panic(err)
	}

	request := anthropic.NewMessageRequest(
		anthropic.WithMessageModel(anthropic.Claude35Sonnet),
		anthropic.WithMessageMaxTokens(20),
		anthropic.WithMessages([]anthropic.MessagePartRequest{{Role: "user", Content: []anthropic.ContentBlock{anthropic.NewTextContentBlock("Hello, world!")}}}),
	)

	response, err := client.Message(ctx, request)
	if err != nil {
		panic(err)
	}

	fmt.Println(response.Content)
}
