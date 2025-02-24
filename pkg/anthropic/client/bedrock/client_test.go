package bedrock

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/pigeonlaser/anthropic-go/v3/pkg/anthropic"
)

func Test_Client_Success_RegionOnly(t *testing.T) {
	client, err := MakeClient(context.Background(), Config{
		Region: "us-west-2",
	})

	assertSuccessClient(t, client, err, "")
}

func Test_Client_Success_RegionWithCredentials(t *testing.T) {
	client, err := MakeClient(context.Background(), Config{
		Region:          "us-west-2",
		AccessKeyID:     "hello-there",
		SecretAccessKey: "general-kenobi",
		SessionToken:    "order-66",
	})

	assertSuccessClient(t, client, err, "")
}

func Test_Client_Success_RegionWithCrossRegionInference(t *testing.T) {
	client, err := MakeClient(context.Background(), Config{
		Region:               "us-west-2",
		CrossRegionInference: true,
	})

	assertSuccessClient(t, client, err, "us")
}

func Test_Client_Failure_MissingRegion(t *testing.T) {
	client, err := MakeClient(context.Background(), Config{})
	if err == nil {
		t.Error("Expected an error when region is not set")
	}

	if client != nil {
		t.Error("Unexpected value for client when region is not set")
	}
}

func Test_Client_Failure_UnsupportedCrossRegionInference(t *testing.T) {
	client, err := MakeClient(context.Background(), Config{
		Region:               "he-llothere",
		CrossRegionInference: true,
	})
	if err == nil {
		t.Error("Expected an error when using an unsupported region for cross region inference")
	}

	if !strings.Contains(err.Error(), "region inference is only supported") {
		t.Errorf("Exepcted an error for unsupported region inference: %s", err.Error())
	}

	if client != nil {
		t.Error("Unexpected value for client when region is not set")
	}
}

type modelTest struct {
	modelInput          anthropic.Model
	expectedModelOutput string
}

func Test_adaptModelForMessage_Success_NonCrossRegion(t *testing.T) {
	client, err := MakeClient(context.Background(), Config{
		Region: "us-west-2",
	})
	if err != nil {
		t.Errorf("Unexpected error when establishing client %s", err.Error())
	}

	testCases := []*modelTest{
		{
			modelInput:          anthropic.Claude35Sonnet,
			expectedModelOutput: BedrockModelClaude35Sonnet,
		},
		{
			modelInput:          anthropic.Claude3Opus,
			expectedModelOutput: BedrockModelClaude3Opus,
		},
		{
			modelInput:          anthropic.Claude3Sonnet,
			expectedModelOutput: BedrockModelClaude3Sonnet,
		},
		{
			modelInput:          anthropic.Claude3Haiku,
			expectedModelOutput: BedrockModelClaude3Haiku,
		},
		{
			modelInput:          anthropic.ClaudeV2_1,
			expectedModelOutput: BedrockModelClaudeV2_1,
		},
	}

	result := ""
	for _, testCase := range testCases {
		result, err = client.adaptModelForMessage(testCase.modelInput)
		if err != nil {
			t.Errorf("Unexpected error when adapting model: %s", err.Error())
		}

		if result != testCase.expectedModelOutput {
			t.Errorf("Error when adapting model. Expected: %s, Actual: %s", testCase.expectedModelOutput, result)
		}
	}
}

func Test_adaptModelForMessage_Failure_UnsupportedModel(t *testing.T) {
	client, err := MakeClient(context.Background(), Config{
		Region: "us-west-2",
	})
	if err != nil {
		t.Errorf("Unexpected error when establishing client %s", err.Error())
	}

	result, err := client.adaptModelForMessage("hello-there")
	if err == nil {
		t.Error("Expected an error when adapting unsupported model")
	}

	if result != "" {
		t.Errorf("Unexpected result for adaptModel: %s", result)
	}
}

func Test_adaptModelForMessage_Success_CrossRegion(t *testing.T) {
	client, err := MakeClient(context.Background(), Config{
		Region:               "eu-west-1",
		CrossRegionInference: true,
	})
	if err != nil {
		t.Errorf("Unexpected error when establishing client %s", err.Error())
	}

	testCases := []*modelTest{
		{
			modelInput:          anthropic.Claude35Sonnet,
			expectedModelOutput: fmt.Sprintf("%s.%s", client.crInferenceRegion, BedrockModelClaude35Sonnet),
		},
		{
			modelInput:          anthropic.Claude3Opus,
			expectedModelOutput: fmt.Sprintf("%s.%s", client.crInferenceRegion, BedrockModelClaude3Opus),
		},
		{
			modelInput:          anthropic.Claude3Sonnet,
			expectedModelOutput: fmt.Sprintf("%s.%s", client.crInferenceRegion, BedrockModelClaude3Sonnet),
		},
		{
			modelInput:          anthropic.Claude3Haiku,
			expectedModelOutput: fmt.Sprintf("%s.%s", client.crInferenceRegion, BedrockModelClaude3Haiku),
		},
	}

	result := ""
	for _, testCase := range testCases {
		result, err = client.adaptModelForMessage(testCase.modelInput)
		if err != nil {
			t.Errorf("Unexpected error when adapting model: %s", err.Error())
		}

		if result != testCase.expectedModelOutput {
			t.Errorf("Error when adapting model. Expected: %s, Actual: %s", testCase.expectedModelOutput, result)
		}
	}
}

func Test_adaptModelForMessage_Failure_ClaudeV2_CrossRegionInference(t *testing.T) {
	client, err := MakeClient(context.Background(), Config{
		Region:               "eu-west-1",
		CrossRegionInference: true,
	})
	if err != nil {
		t.Errorf("Unexpected error when establishing client %s", err.Error())
	}

	result, err := client.adaptModelForMessage(anthropic.ClaudeV2_1)
	if err == nil {
		t.Error("Expected an error when using cross region inference on claude v2.1")
	}

	if !strings.Contains(err.Error(), "not compatible with cross-region") {
		t.Error("Expected a 'not compatible with cross-region' error")
	}

	if result != "" {
		t.Errorf("Unexpected result for adaptModel: %s", result)
	}
}

func assertSuccessClient(t *testing.T, client *Client, err error, crRegionValue string) {
	if err != nil {
		t.Errorf("Unexpected error %s", err.Error())
	}

	if client.brCli == nil {
		t.Error("Unexpected nil for brCli")
	}

	if client.crInferenceRegion != crRegionValue {
		t.Errorf("Unexpected value for inference region %s", client.crInferenceRegion)
	}
}
