package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/madebywelch/anthropic-go/v2/pkg/anthropic"
)

func main() {
	client, err := anthropic.NewClient("your-api-key")
	if err != nil {
		panic(err)
	}

	messages := []anthropic.MessagePartRequest{
		{
			Role:    "user",
			Content: "Why is the sky blue?",
		},
	}

	request := anthropic.NewMessageRequest(
		messages,
		anthropic.WithModel[anthropic.MessageRequest](anthropic.Claude3Opus),
		anthropic.WithMaxTokens[anthropic.MessageRequest](1000),
		anthropic.WithStream[anthropic.MessageRequest](true),
	)

	rCh, errCh := client.MessageStream(request)

	final := strings.Builder{}
	chunk := anthropic.MessageStreamResponse{}
	done := false

	for {
		select {
		case chunk = <-rCh:
			final.WriteString(chunk.Delta.Text)
			fmt.Print(chunk.Delta.Text)
		case err := <-errCh:
			fmt.Printf("\n\nError: %s\n\n", err)
			done = true
			break
		}
		if chunk.Type == "message_stop" || done {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("-------------------FINAL RESULT----------------------")
	fmt.Println(final.String())
	fmt.Println("-----------------------------------------------------")
}
