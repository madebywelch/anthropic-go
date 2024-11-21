package bedrock

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/madebywelch/anthropic-go/v3/pkg/anthropic"
)

func (c *Client) Complete(ctx context.Context, req *anthropic.CompletionRequest) (*anthropic.CompletionResponse, error) {
	err := anthropic.ValidateCompleteRequest(req)
	if err != nil {
		return nil, err
	}

	return c.sendCompleteRequest(ctx, req)
}

func (c *Client) sendCompleteRequest(ctx context.Context, req *anthropic.CompletionRequest) (*anthropic.CompletionResponse, error) {
	adaptedModel, err := adaptModelForCompletion(req.Model)
	if err != nil {
		return nil, err
	}

	// Adapt the request to a Bedrock request
	bedReq := adaptCompletionRequest(req)

	data, err := json.Marshal(bedReq)
	if err != nil {
		return nil, fmt.Errorf("error marshalling complete request: %w", err)
	}

	response, err := c.brCli.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		Body:        data,
		ModelId:     aws.String(adaptedModel),
		ContentType: aws.String("application/json"),
	})

	if err != nil {
		errStatusCode := extractErrStatusCode(err)
		return nil, anthropic.MapHTTPStatusCodeToError(errStatusCode, err.Error())
	}

	compResp := &anthropic.CompletionResponse{}
	err = json.Unmarshal(response.Body, compResp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling complete response: %w", err)
	}

	return compResp, nil
}
