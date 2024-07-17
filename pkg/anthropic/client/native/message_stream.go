package native

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/madebywelch/anthropic-go/v3/pkg/anthropic"
)

func (c *Client) MessageStream(ctx context.Context, req *anthropic.MessageRequest) (<-chan *anthropic.MessageStreamResponse, <-chan error) {
	msCh := make(chan *anthropic.MessageStreamResponse)
	errCh := make(chan error, 1)

	err := anthropic.ValidateMessageStreamRequest(req)
	if err != nil {
		errCh <- err
		close(msCh)
		close(errCh)
		return msCh, errCh
	}

	go c.handleMessageStreaming(ctx, req, msCh, errCh)

	return msCh, errCh
}

func (c *Client) handleMessageStreaming(
	ctx context.Context,
	req *anthropic.MessageRequest,
	msCh chan<- *anthropic.MessageStreamResponse,
	errCh chan<- error,
) {
	defer close(msCh)
	defer close(errCh)

	data, err := json.Marshal(req)
	if err != nil {
		errCh <- fmt.Errorf("error marshalling message request: %w", err)
		return
	}

	request, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/v1/messages", c.baseURL), bytes.NewBuffer(data))
	if err != nil {
		errCh <- fmt.Errorf("error creating new request: %w", err)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Api-Key", c.apiKey)
	request.Header.Set("Accept", "text/event-stream")
	if len(req.Beta) > 0 {
		request.Header.Set("anthropic-beta", req.Beta)
	}

	response, err := c.doRequest(request)
	if err != nil {
		errCh <- fmt.Errorf("error sending message request: %w", err)
		return
	}
	defer response.Body.Close()

	err = c.processMessageSseStream(response.Body, msCh)
	if err != nil {
		errCh <- err
	}
}

func (c *Client) processMessageSseStream(reader io.Reader, events chan<- *anthropic.MessageStreamResponse) error {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "data:") {
			data := strings.TrimSpace(line[5:])

			event := &anthropic.MessageEvent{}
			err := json.Unmarshal([]byte(data), event)
			if err != nil {
				return fmt.Errorf("error decoding event data: %w", err)
			}

			msg, err := anthropic.ParseMessageEvent(anthropic.MessageEventType(event.Type), data)

			if err != nil {
				if _, ok := err.(anthropic.UnsupportedEventType); ok {
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
