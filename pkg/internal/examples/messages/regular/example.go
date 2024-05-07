package main

/*
func main() {
	client, err := anthropic.NewClient("your-api-key")
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
	response, err := client.Message(request)
	if err != nil {
		panic(err)
	}

	fmt.Println(response.Content)
}
*/
