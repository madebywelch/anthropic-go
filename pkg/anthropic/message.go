package anthropic

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (c *Client) Message(req *MessageRequest) (*MessageResponse, error) {
	if req.Stream {
		return nil, fmt.Errorf("cannot use Message with streaming enabled, use MessageStream instead")
	}

	if !req.Model.IsMessageCompatible() {
		return nil, fmt.Errorf("model %s is not compatible with the message endpoint", req.Model)
	}

	if !req.Model.IsImageCompatible() && req.ContainsImageContent() {
		return nil, fmt.Errorf("model %s does not support image content", req.Model)
	}

	if req.CountImageContent() > 20 {
		return nil, fmt.Errorf("too many image content blocks, maximum is 20")
	}

	return c.sendMessageRequest(req)
}

func (c *Client) MessageStream(req *MessageRequest) (<-chan MessageStreamResponse, <-chan error) {
	events := make(chan MessageStreamResponse)

	// make this a buffered channel to allow for the error case below to return
	errCh := make(chan error, 1)

	if !req.Stream {
		errCh <- fmt.Errorf("cannot use MessageStream with a non-streaming request, use Message instead")
		return events, errCh
	}

	if !req.Model.IsMessageCompatible() {
		errCh <- fmt.Errorf("model %s is not compatible with the messagestream endpoint", req.Model)
		return events, errCh
	}

	if req.Stream && len(req.Tools) > 0 {
		// https://docs.anthropic.com/claude/docs/tool-use
		// Streaming (stream=true) is not yet supported. We plan to add streaming support in a future beta version.
		errCh <- fmt.Errorf("cannot use streaming with tools")
		return events, errCh
	}

	go c.handleMessageStreaming(events, errCh, req)

	return events, errCh
}

func (c *Client) sendMessageRequest(req *MessageRequest) (*MessageResponse, error) {
	// Marshal the request to JSON
	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshalling completion request: %w", err)
	}

	// Create the HTTP request
	requestURL := fmt.Sprintf("%s/v1/messages", c.baseURL)
	request, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("error creating new request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Api-Key", c.apiKey)
	request.Header.Set("anthropic-beta", AnthropicAPIToolsBeta)

	// Use the DoRequest method to send the HTTP request
	response, err := c.doRequest(request)
	if err != nil {
		return nil, fmt.Errorf("error sending completion request: %w", err)
	}
	defer response.Body.Close()

	// Decode the response body to a MessageResponse object
	var messageResponse MessageResponse
	err = json.NewDecoder(response.Body).Decode(&messageResponse)
	if err != nil {
		return nil, fmt.Errorf("error decoding message response: %w", err)
	}

	return &messageResponse, nil
}

func (c *Client) handleMessageStreaming(events chan MessageStreamResponse, errCh chan error, req *MessageRequest) {
	defer close(events)

	data, err := json.Marshal(req)
	if err != nil {
		errCh <- fmt.Errorf("error marshalling message request: %w", err)
		return
	}

	request, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/messages", c.baseURL), bytes.NewBuffer(data))
	if err != nil {
		errCh <- fmt.Errorf("error creating new request: %w", err)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Api-Key", c.apiKey)
	request.Header.Set("Accept", "text/event-stream")

	response, err := c.doRequest(request)
	if err != nil {
		errCh <- fmt.Errorf("error sending message request: %w", err)
		return
	}
	defer response.Body.Close()

	err = c.processMessageSseStream(response.Body, events)
	if err != nil {
		errCh <- err
	}
}

func (c *Client) processMessageSseStream(reader io.Reader, events chan MessageStreamResponse) error {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "data:") {
			data := strings.TrimSpace(line[5:])

			event := &MessageEvent{}
			err := json.Unmarshal([]byte(data), event)
			if err != nil {
				return fmt.Errorf("error decoding event data: %w", err)
			}

			msg, err := parseMessageEvent(event.Type, data)

			if err != nil {
				if _, ok := err.(UnsupportedEventType); ok {
					// ignore unsupported event types
				} else {
					return fmt.Errorf("error processing message stream: %v", err)
				}
			}

			events <- msg
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading from stream: %w", err)
	}

	return nil
}
