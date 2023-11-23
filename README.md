# Unofficial Anthropic SDK in Go

This project provides an unofficial Go SDK for Anthropic, a A next-generation AI assistant for your tasks, no matter the scale. The SDK makes it easy to interact with the Anthropic API in Go applications. For more information about Anthropic, including API documentation, visit the official [Anthropic documentation.](https://console.anthropic.com/docs)

[![GoDoc](https://godoc.org/github.com/madebywelch/anthropic-go?status.svg)](https://pkg.go.dev/github.com/madebywelch/anthropic-go/v2)

## Installation

You can install the Anthropic SDK in Go using go get:

```go
go get github.com/madebywelch/anthropic-go/v2
```

## Usage

To use the Anthropic SDK, you'll need to initialize a client and make requests to the Anthropic API. Here's an example of initializing a client and performing a regular and a streaming completion:

## Completion Example

```go
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
		anthropic.WithModel(anthropic.ClaudeV2_1),
		anthropic.WithMaxTokens(100),
	)

	// Note: Only use client.Complete when streaming is disabled, otherwise use client.CompleteStream!
	response, err := client.Complete(request)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Completion: %s\n", response.Completion)
}
```

### Completion Example Output

```
The sky appears blue to us due to the way the atmosphere scatters light from the sun
```

## Streaming Example

```go
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
		anthropic.WithModel(anthropic.ClaudeV2_1),
		anthropic.WithMaxTokens(100),
		anthropic.WithStreaming(true),
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
```

### Streaming Example Output

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

## Contributing

Contributions to this project are welcome. To contribute, follow these steps:

- Fork this repository
- Create a new branch (`git checkout -b feature/my-new-feature`)
- Commit your changes (`git commit -am 'Add some feature'`)
- Push the branch (`git push origin feature/my-new-feature`)
- Create a new pull request

## License

This project is licensed under the Apache License, Version 2.0 - see the [LICENSE](LICENSE) file for details.
