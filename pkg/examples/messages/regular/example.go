package main

import (
	"fmt"

	"github.com/madebywelch/anthropic-go/v2/pkg/anthropic"
)

func main() {
	client, err := anthropic.NewClient("your-api-key")
	if err != nil {
		panic(err)
	}

	// Prepare a message request
	request := &anthropic.MessageRequest{
		Model:             anthropic.ClaudeV2_1,
		MaxTokensToSample: 10,
		Messages:          []anthropic.MessagePartRequest{{Role: "user", Content: "Hello, Anthropics!"}},
	}

	// Call the Message method
	response, err := client.Message(request)
	if err != nil {
		panic(err)
	}

	fmt.Println(response.Content)
}
