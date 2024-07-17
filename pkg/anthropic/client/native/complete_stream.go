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

func (c *Client) CompleteStream(ctx context.Context, req *anthropic.CompletionRequest) (<-chan *anthropic.StreamResponse, <-chan error) {
	cCh := make(chan *anthropic.StreamResponse)
	errCh := make(chan error, 1)

	err := anthropic.ValidateCompleteStreamRequest(req)
	if err != nil {
		errCh <- err
		close(cCh)
		close(errCh)
		return cCh, errCh
	}

	go c.handleCompleteStreaming(ctx, req, cCh, errCh)
	return cCh, errCh
}

func (c *Client) handleCompleteStreaming(
	ctx context.Context,
	req *anthropic.CompletionRequest,
	cCh chan<- *anthropic.StreamResponse,
	errCh chan<- error,
) {
	defer close(cCh)
	defer close(errCh)

	data, err := json.Marshal(req)
	if err != nil {
		errCh <- fmt.Errorf("error marshalling completion request: %w", err)
		return
	}

	request, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/v1/complete", c.baseURL), bytes.NewBuffer(data))
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
		errCh <- fmt.Errorf("error sending completion request: %w", err)
		return
	}
	defer response.Body.Close()

	err = c.processSseStream(response.Body, cCh)
	if err != nil {
		errCh <- err
	}
}

func (c *Client) processSseStream(reader io.Reader, cCh chan<- *anthropic.StreamResponse) error {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "data:") {
			data := strings.TrimSpace(line[5:])
			event := &anthropic.StreamResponse{}
			err := json.Unmarshal([]byte(data), event)
			if err != nil {
				return fmt.Errorf("error decoding event data: %w", err)
			}

			cCh <- event
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading from stream: %w", err)
	}

	return nil
}
