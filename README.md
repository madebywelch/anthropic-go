# Unofficial Anthropic SDK for Go

This project provides an unofficial Go SDK for Anthropic, a next-generation AI assistant platform. The SDK simplifies interactions with the Anthropic API in Go applications. For more information about Anthropic and its API, visit the [official Anthropic documentation](https://console.anthropic.com/docs).

[![GoDoc](https://godoc.org/github.com/madebywelch/anthropic-go?status.svg)](https://pkg.go.dev/github.com/madebywelch/anthropic-go/v4)

## Installation

Install the Anthropic SDK for Go using:

```go
go get github.com/madebywelch/anthropic-go/v4
```

## Features

- Support for both native Anthropic API and AWS Bedrock
- Completion and streaming completion
- Message and streaming message support
- Tool usage capabilities

## Quick Start

Here's a basic example of using the SDK:

```go
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
		[]anthropic.MessagePartRequest{{Role: "user", Content: []anthropic.ContentBlock{anthropic.NewTextContentBlock("Hello, world!")}}},
		anthropic.WithModel[anthropic.MessageRequest](anthropic.Claude35Sonnet),
		anthropic.WithMaxTokens[anthropic.MessageRequest](20),
	)

	response, err := client.Message(ctx, request)
	if err != nil {
		panic(err)
	}

	fmt.Println(response.Content)
}
```

## Usage

For more detailed usage examples, including streaming, completion, and tool usage, please refer to the `pkg/internal/examples` directory in the repository.

## Contributing

Contributions are welcome! To contribute:

1. Fork this repository
2. Create a new branch (`git checkout -b feature/my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push the branch (`git push origin feature/my-new-feature`)
5. Create a new pull request

## License

This project is licensed under the Apache License, Version 2.0. See the [LICENSE](LICENSE) file for details.
