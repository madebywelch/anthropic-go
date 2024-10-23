package bedrock

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/madebywelch/anthropic-go/v3/pkg/anthropic"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

const (
	AnthropicVersion = "bedrock-2023-05-31"

	BedrockModelClaude35Sonnet_20241022 = "anthropic.claude-3-5-sonnet-20241022-v2:0"
	BedrockModelClaude35Sonnet          = "anthropic.claude-3-5-sonnet-20240620-v1:0"
	BedrockModelClaude3Opus             = "anthropic.claude-3-opus-20240229-v1:0"
	BedrockModelClaude3Sonnet           = "anthropic.claude-3-sonnet-20240229-v1:0"
	BedrockModelClaude3Haiku            = "anthropic.claude-3-haiku-20240307-v1:0"
	BedrockModelClaudeV2_1              = "anthropic.claude-v2:1"

	// Cross-region top-level region code
	CRUS = "us"
	CREU = "eu"
)

type Client struct {
	brCli             *bedrockruntime.Client
	crInferenceRegion string
}

type Config struct {
	Region               string
	AccessKeyID          string
	SecretAccessKey      string
	SessionToken         string
	CrossRegionInference bool
}

func MakeClient(ctx context.Context, cfg Config) (*Client, error) {
	if cfg.Region == "" {
		return nil, fmt.Errorf("Region is requried for establishing anthropic bedrock client")
	}

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

	regionPrefix := ""
	if cfg.CrossRegionInference {
		// extract the first 2 letters from the region
		regionPrefix = cfg.Region[:2]
		if regionPrefix != CRUS && regionPrefix != CREU {
			return nil, fmt.Errorf(
				"Cross region inference is only supported for: '%s', '%s'; Region prefix: '%s' is not supported",
				CRUS,
				CREU,
				regionPrefix,
			)
		}
	}

	return &Client{
		brCli:             bedrockruntime.NewFromConfig(awsCfg),
		crInferenceRegion: regionPrefix,
	}, nil
}

// adaptModelForMessage takes the model as defined in anthropic.Model and adapts it to the model Bedrock expects
func (c *Client) adaptModelForMessage(model anthropic.Model) (string, error) {
	adaptedModel := ""

	switch model {
	case anthropic.Claude35_Sonnet_20241022:
		adaptedModel = BedrockModelClaude35Sonnet_20241022
	case anthropic.Claude35Sonnet:
		adaptedModel = BedrockModelClaude35Sonnet
	case anthropic.Claude3Opus:
		adaptedModel = BedrockModelClaude3Opus
	case anthropic.Claude3Sonnet:
		adaptedModel = BedrockModelClaude3Sonnet
	case anthropic.Claude3Haiku:
		adaptedModel = BedrockModelClaude3Haiku
	case anthropic.ClaudeV2_1:
		adaptedModel = BedrockModelClaudeV2_1
	default:
		return "", fmt.Errorf("model %s is not compatible with the bedrock message endpoint", model)
	}

	if c.crInferenceRegion == "" {
		return adaptedModel, nil
	}

	if adaptedModel == BedrockModelClaudeV2_1 {
		return "", fmt.Errorf("Bedrock model %s is not compatible with cross-region inference", adaptedModel)
	}

	return fmt.Sprintf("%s.%s", c.crInferenceRegion, adaptedModel), nil
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
