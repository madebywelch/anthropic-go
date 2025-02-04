package bedrock

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/madebywelch/anthropic-go/v4/pkg/anthropic"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
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

	adaptedModel, err := c.adaptModelForMessage(req.Model)
	if err != nil {
		errCh <- fmt.Errorf("error adapting model: %w", err)
		return
	}

	// Adapt the request to a Bedrock request
	bedReq := adaptMessageRequest(req)

	data, err := json.Marshal(bedReq)
	if err != nil {
		errCh <- fmt.Errorf("error marshalling message request: %w", err)
		return
	}

	response, err := c.brCli.InvokeModelWithResponseStream(
		ctx,
		&bedrockruntime.InvokeModelWithResponseStreamInput{
			Body:        data,
			ModelId:     aws.String(adaptedModel),
			ContentType: aws.String("application/json"),
		},
	)
	if err != nil {
		errStatusCode := extractErrStatusCode(err)
		errCh <- anthropic.MapHTTPStatusCodeToError(errStatusCode)
		return
	}

	for event := range response.GetStream().Events() {
		select {
		case <-ctx.Done():
			return
		default:
		}

		if v, ok := event.(*types.ResponseStreamMemberChunk); ok {
			event := &anthropic.MessageEvent{}
			err := json.Unmarshal(v.Value.Bytes, event)
			if err != nil {
				errCh <- fmt.Errorf("error decoding event data: %w", err)
				return
			}
			msg, err := anthropic.ParseMessageEvent(
				anthropic.MessageEventType(event.Type),
				string(v.Value.Bytes),
			)
			if err != nil {
				if _, ok := err.(anthropic.UnsupportedEventType); ok {
					// ignore unsupported event types
				} else {
					errCh <- fmt.Errorf("error processing message stream: %v", err)
					return
				}
			}

			msCh <- msg
		}
	}
}
