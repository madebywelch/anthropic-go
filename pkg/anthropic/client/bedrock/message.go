package bedrock

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/madebywelch/anthropic-go/v3/pkg/anthropic"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

func (c *Client) Message(ctx context.Context, req *anthropic.MessageRequest) (*anthropic.MessageResponse, error) {
	err := anthropic.ValidateMessageRequest(req)
	if err != nil {
		return nil, err
	}

	return c.sendMessageRequest(ctx, req)
}

func (c *Client) sendMessageRequest(ctx context.Context, req *anthropic.MessageRequest) (*anthropic.MessageResponse, error) {
	adaptedModel, err := c.adaptModelForMessage(req.Model)
	if err != nil {
		return nil, err
	}

	// Adapt the request to a Bedrock request
	bedReq := adaptMessageRequest(req)

	data, err := json.Marshal(bedReq)
	if err != nil {
		return nil, fmt.Errorf("error marshalling message request: %w", err)
	}

	response, err := c.brCli.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		Body:        data,
		ModelId:     aws.String(adaptedModel),
		ContentType: aws.String("application/json"),
	})

	if err != nil {
		errStatusCode := extractErrStatusCode(err)
		return nil, anthropic.MapHTTPStatusCodeToError(errStatusCode)
	}

	msgResp := &anthropic.MessageResponse{}
	err = json.Unmarshal(response.Body, msgResp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling message response: %w", err)
	}

	return msgResp, nil
}
