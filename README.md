# Unofficial Anthropic SDK in Go

This project provides an unofficial Go SDK for Anthropic, a A next-generation AI assistant for your tasks, no matter the scale. The SDK makes it easy to interact with the Anthropic API in Go applications. For more information about Anthropic, including API documentation, visit the official [Anthropic documentation.](https://console.anthropic.com/docs)

[![GoDoc](https://godoc.org/github.com/madebywelch/anthropic-go?status.svg)](https://pkg.go.dev/github.com/madebywelch/anthropic-go/v2)

## Installation

You can install the Anthropic SDK in Go using go get:

```go
go get github.com/madebywelch/anthropic-go/v2
```

## Client Instantiation

The native Anthropic API and access to Anthropic via AWS Bedrock are both supported. Here's an example of instantiating each client:

## Client Instantiation Examples

```go
package main

import (
    "context"

    "github.com/madebywelch/anthropic-go/v3/pkg/anthropic/client"
    "github.com/madebywelch/anthropic-go/v3/pkg/anthropic/client/native"
    "github.com/madebywelch/anthropic-go/v3/pkg/anthropic/client/bedrock"
)


func main() {
    brCli := client.MakeClient(ctx, &bedrock.Config{
        Region: "us-west-2", // us-west-2 is the only region that supports opus as of 2024-05-09
        AccessKeyID: "", // Optional, only required if IAM role isn't set up
        SecretAccessKey: "", // Optional, only required if IAM role isn't set up
        SessionToken: "", // Optional, only required if IAM role isn't set up
    })

    nativeCli := client.MakeClient(ctx, &native.Config{
        APIKey: "your-anthropic-api-key",
    })

    /*
        both clients implement the same interface and have access to:
        - Complete
        - CompleteStream
        - Message
        - MessageStream
    */
}
```

## Usage

To use the Anthropic SDK, you'll need to initialize a client and make requests to the Anthropic API. Here's an example of initializing a client and performing a regular and a streaming completion:

## Completion Example

```go
package main

import (
	"context"
	"fmt"

	"github.com/madebywelch/anthropic-go/v3/pkg/anthropic"
	"github.com/madebywelch/anthropic-go/v3/pkg/anthropic/client"
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
		anthropic.WithStreaming[anthropic.CompletionRequest](true),
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
```

### Completion Example Output

```
The sky appears blue to us due to the way the atmosphere scatters light from the sun
```

## Completion Streaming Example

```go
package main

import (
	"context"
	"fmt"

	"github.com/madebywelch/anthropic-go/v3/pkg/anthropic"
	"github.com/madebywelch/anthropic-go/v3/pkg/anthropic/client"
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
		anthropic.WithStreaming[anthropic.CompletionRequest](true),
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
```

### Completion Streaming Example Output

```
There
are
a
few
reasons
why
the
sky
appears
```

## Messages Example

```go
package main

import (
	"context"
	"fmt"

	"github.com/madebywelch/anthropic-go/v3/pkg/anthropic"
	"github.com/madebywelch/anthropic-go/v3/pkg/anthropic/client"
	"github.com/madebywelch/anthropic-go/v3/pkg/anthropic/client/native"
)

func main() {
	ctx := context.Background()
	client, err := client.MakeClient(ctx, &native.Config{
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
```

### Messages Example Output

```
{ID:msg_01W3bZkuMrS3h1ehqTdF84vv Type:message Model:claude-2.1 Role:assistant Content:[{Type:text Text:Hello!}] StopReason:end_turn Stop: StopSequence:}
```

## Messages Streaming Example

```go
package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/madebywelch/anthropic-go/v3/pkg/anthropic"
	"github.com/madebywelch/anthropic-go/v3/pkg/anthropic/client"
	"github.com/madebywelch/anthropic-go/v3/pkg/anthropic/client/native"
)

func main() {
	ctx := context.Background()
	client, err := client.MakeClient(ctx, &native.Config{
		APIKey: "your-api-key",
	})
	if err != nil {
		panic(err)
	}

	messages := []anthropic.MessagePartRequest{
		{
			Role:    "user",
			Content: []anthropic.ContentBlock{anthropic.NewTextContentBlock("Hello, world!")},
		},
	}

	request := anthropic.NewMessageRequest(
		messages,
		anthropic.WithModel[anthropic.MessageRequest](anthropic.Claude3Opus),
		anthropic.WithMaxTokens[anthropic.MessageRequest](1000),
		anthropic.WithStream[anthropic.MessageRequest](true),
	)

	rCh, errCh := client.MessageStream(ctx, request)

	final := strings.Builder{}
	chunk := &anthropic.MessageStreamResponse{}
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
}
```

### Messages Streaming Example Output

```
Good
 morning
!
 As
 an
 AI
 language
 model
,
 I
 don
't
 have
 feelings
 or
 a
 physical
 state
,
 but
```

## Messages Tools Example

```go
package main

import (
	"context"

	"github.com/madebywelch/anthropic-go/v3/pkg/anthropic"
	"github.com/madebywelch/anthropic-go/v3/pkg/anthropic/client"
	"github.com/madebywelch/anthropic-go/v3/pkg/anthropic/client/native"
)

func main() {
	ctx := context.Background()
	client, err := client.MakeClient(ctx, &native.Config{
		APIKey: "your-api-key",
	})
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
	response, err := client.Message(ctx, request)
	if err != nil {
		panic(err)
	}

	if response.StopReason == "tool_use" {
		// Do something with the tool response
	}
}
```

## Contributing

Contributions to this project are welcome. To contribute, follow these steps:

- Fork this repository
- Create a new branch (`git checkout -b feature/my-new-feature`)
- Commit your changes (`git commit -am 'Add some feature'`)
- Push the branch (`git push origin feature/my-new-feature`)
- Create a new pull request

## License

This project is licensed under the Apache License, Version 2.0 - see the [LICENSE](LICENSE) file for details.
