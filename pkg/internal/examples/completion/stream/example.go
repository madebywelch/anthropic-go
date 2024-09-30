package main

import (
	"context"
	"fmt"

	"github.com/madebywelch/anthropic-go/v4/pkg/anthropic"
	"github.com/madebywelch/anthropic-go/v4/pkg/anthropic/client/native"
	"github.com/madebywelch/anthropic-go/v4/pkg/anthropic/utils"
)

func main() {
	ctx := context.Background()
	client, err := native.MakeClient(native.Config{
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
		anthropic.WithCompletionModel(anthropic.ClaudeV2_1),
		anthropic.WithCompletionMaxTokens(100),
		anthropic.WithCompletionStream(true),
	)

	// Note: Only use client.CompleteStream when streaming is enabled, otherwise use client.Complete!
	resps, errs := client.CompleteStream(ctx, request)

	for {
		select {
		case resp := <-resps:
			fmt.Printf("Completion: %s\n", resp.Completion)
		case err := <-errs:
			panic(err)
		}
	}
}
