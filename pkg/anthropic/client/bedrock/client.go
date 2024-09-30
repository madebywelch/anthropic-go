package bedrock

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/madebywelch/anthropic-go/v4/pkg/anthropic"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

const (
	AnthropicVersion = "bedrock-2023-05-31"

	BedrockModelClaude3Opus   = "anthropic.claude-3-opus-20240229-v1:0"
	BedrockModelClaude3Sonnet = "anthropic.claude-3-sonnet-20240229-v1:0"
	BedrockModelClaude3Haiku  = "anthropic.claude-3-haiku-20240307-v1:0"
	BedrockModelClaudeV2_1    = "anthropic.claude-v2:1"
)

type Client struct {
	brCli *bedrockruntime.Client
}

type Config struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
}

func MakeClient(ctx context.Context, cfg Config) (*Client, error) {
	awsCfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(cfg.Region),
	)

	// override config load with static credentials if provided
	if cfg.AccessKeyID != "" && cfg.SecretAccessKey != "" {
		credsProvider := credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, cfg.SessionToken)
		awsCfg, err = config.LoadDefaultConfig(
			ctx,
			config.WithRegion(cfg.Region),
			config.WithCredentialsProvider(credsProvider),
		)
	}

	if err != nil {
		return nil, err
	}

	return &Client{
		brCli: bedrockruntime.NewFromConfig(awsCfg),
	}, nil
}

// adaptModelForMessage takes the model as defined in anthropic.Model and adapts it to the model Bedrock expects
func adaptModelForMessage(model anthropic.Model) (string, error) {
	if model == anthropic.Claude3Opus {
		return BedrockModelClaude3Opus, nil
	}
	if model == anthropic.Claude3Sonnet {
		return BedrockModelClaude3Sonnet, nil
	}
	if model == anthropic.Claude3Haiku {
		return BedrockModelClaude3Haiku, nil
	}
	if model == anthropic.ClaudeV2_1 {
		return BedrockModelClaudeV2_1, nil
	}

	return "", fmt.Errorf("model %s is not compatible with the bedrock message endpoint", model)
}

// adaptModelForCompletion takes the model as defined in anthropic.Model and adapts it to the model Bedrock expects
func adaptModelForCompletion(model anthropic.Model) (string, error) {
	if model == anthropic.ClaudeV2_1 {
		return BedrockModelClaudeV2_1, nil
	}

	return "", fmt.Errorf("model %s is not compatible with the bedrock completion endpoint", model)
}

// MessageRequest is an override for the default message request to adapt the request for the Bedrock API.
type MessageRequest struct {
	anthropic.MessageRequest
	AnthropicVersion string `json:"anthropic_version"`
	Model            bool   `json:"model,omitempty"`  // shadow for Model
	Stream           bool   `json:"stream,omitempty"` // shadow for Stream
}

func adaptMessageRequest(req *anthropic.MessageRequest) *MessageRequest {
	return &MessageRequest{
		MessageRequest:   *req,
		AnthropicVersion: AnthropicVersion,
	}
}

type CompleteRequest struct {
	anthropic.CompletionRequest
	AnthropicVersion string `json:"anthropic_version"`
	Model            bool   `json:"model,omitempty"`  // shadow for Model
	Stream           bool   `json:"stream,omitempty"` // shadow for Stream
}

func adaptCompletionRequest(req *anthropic.CompletionRequest) *CompleteRequest {
	return &CompleteRequest{
		CompletionRequest: *req,
		AnthropicVersion:  AnthropicVersion,
	}
}

func extractErrStatusCode(err error) int {
	re := regexp.MustCompile(`StatusCode: (\d+)`)
	match := re.FindStringSubmatch(err.Error())

	if len(match) > 1 {
		res, err := strconv.Atoi(match[1])
		if err != nil {
			return 0
		}
		return res
	}

	return 0
}
