package main

import (
	"fmt"

	"github.com/madebywelch/anthropic-go/pkg/anthropic"
	"github.com/madebywelch/anthropic-go/pkg/anthropic/utils"
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
		anthropic.WithModel(anthropic.ClaudeV1),
		anthropic.WithMaxTokens(100),
	)

	// Note: Only use client.Complete when streaming is disabled, otherwise use client.CompleteStream!
	response, err := client.Complete(request)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Completion: %s\n", response.Completion)
}
