package main

import (
	"fmt"

	"github.com/madebywelch/anthropic-go/v2/pkg/anthropic"
	"github.com/madebywelch/anthropic-go/v2/pkg/anthropic/utils"
)

func main() {
	client, err := anthropic.NewClient("your-api-key")
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
		anthropic.WithStreaming[anthropic.CompletionRequest](true),
	)

	// Note: Only use client.CompleteStream when streaming is enabled, otherwise use client.Complete!
	resps, errs := client.CompleteStream(request)

	for {
		select {
		case resp := <-resps:
			fmt.Printf("Completion: %s\n", resp.Completion)
		case err := <-errs:
			panic(err)
		}
	}
}
