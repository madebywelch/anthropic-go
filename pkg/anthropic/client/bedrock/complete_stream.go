package bedrock

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/madebywelch/anthropic-go/v3/pkg/anthropic"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
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

	adaptedModel, err := adaptModelForCompletion(req.Model)
	if err != nil {
		errCh <- err
		return
	}

	// Adapt the request to a Bedrock request
	bedReq := adaptCompletionRequest(req)

	data, err := json.Marshal(bedReq)
	if err != nil {
		errCh <- fmt.Errorf("error marshalling complete request: %w", err)
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
		errCh <- err
		return
	}

	for event := range response.GetStream().Events() {
		select {
		case <-ctx.Done():
			return
		default:
		}

		if v, ok := event.(*types.ResponseStreamMemberChunk); ok {
			streamResp := &anthropic.StreamResponse{}
			err = json.Unmarshal(v.Value.Bytes, streamResp)
			if err != nil {
				errCh <- fmt.Errorf("error unmarshalling stream response: %w", err)
				return
			}

			fmt.Println(streamResp)

			cCh <- streamResp
		}
	}
}
