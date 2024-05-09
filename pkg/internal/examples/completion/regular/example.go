package main

import (
	"context"
	"fmt"

	anthropic "github.com/madebywelch/anthropic-go/v3/pkg/anthropic"
	client "github.com/madebywelch/anthropic-go/v3/pkg/anthropic/client"
	"github.com/madebywelch/anthropic-go/v3/pkg/anthropic/client/native"
	"github.com/madebywelch/anthropic-go/v3/pkg/anthropic/utils"
)

func main() {
	ctx := context.Background()
	client, err := client.MakeClient(ctx, &native.Config{
		APIKey: "your-api-key",
	})
	if err != nil {
		panic(err)
	}

	prompt, err := utils.GetPrompt("Why is the sky blue?")
	if err != nil {
		panic(err)
	}

	request := anthropic.NewCompletionRequest(
		prompt,
		anthropic.WithModel[anthropic.CompletionRequest](anthropic.ClaudeV2_1),
		anthropic.WithMaxTokens[anthropic.CompletionRequest](100),
	)

	// Note: Only use client.Complete when streaming is disabled, otherwise use client.CompleteStream!
	response, err := client.Complete(ctx, request)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Completion: %s\n", response.Completion)
}
