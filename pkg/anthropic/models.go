package anthropic

// Model represents a Claude model.
type Model string

// DefaultModel is the default model used if none is specified.
const DefaultModel = Claude3Sonnet

// https://docs.anthropic.com/claude/docs/models-overview
const (
	// Highest level of intelligence and capability
	Claude35Sonnet Model = "claude-3-5-sonnet-latest"

	// New version of claude-3-5-sonnet
	Claude35Sonnet_20241022 Model = "claude-3-5-sonnet-20241022"

	// Former highest level of intelligence and capability
	Claude35Sonnet_20240620 Model = "claude-3-5-sonnet-20240620"

	// Fastest and most compact model for near-instant responsiveness
	Claude35Haiku Model = "claude-3-5-haiku-latest"

	// Fastest and most compact model for near-instant responsiveness, 20241022 model
	Claude35Haiku_20241022 Model = "claude-3-5-haiku-20241022"

	// Updated version of Claude 2 with improved accuracy
	ClaudeV2_1 Model = "claude-2.1"
	ClaudeV2   Model = "claude-2"

	// Claude 1 models
	ClaudeV1        Model = "claude-v1"
	ClaudeV1_100k   Model = "claude-v1-100k"
	ClaudeV1_3      Model = "claude-v1.3"
	ClaudeV1_3_100k Model = "claude-v1.3-100k"
	ClaudeV1_2      Model = "claude-v1.2"
	ClaudeV1_0      Model = "claude-v1.0"

	// Claude Instant models
	ClaudeInstantV1        Model = "claude-instant-v1"
	ClaudeInstantV1_100k   Model = "claude-instant-v1-100k"
	ClaudeInstantV1_1      Model = "claude-instant-v1.1"
	ClaudeInstantV1_1_100k Model = "claude-instant-v1.1-100k"
	ClaudeInstantV1_0      Model = "claude-instant-v1.0"

	// First, add the missing Claude 3 model constants
	Claude3Opus   Model = "claude-3-opus-20240229"
	Claude3Sonnet Model = "claude-3-sonnet-20240229"
	Claude3Haiku  Model = "claude-3-haiku-20240307"
)

// Define the imageCompatibleModels map alongside the other maps
var (
	imageCompatibleModels = map[Model]bool{
		Claude3Opus:             true,
		Claude3Sonnet:           true,
		Claude3Haiku:            true,
		Claude35Sonnet:          true,
		Claude35Sonnet_20241022: true,
		Claude35Sonnet_20240620: true,
	}

	messageCompatibleModels = map[Model]bool{
		Claude35Sonnet:          true,
		Claude35Sonnet_20241022: true,
		Claude35Sonnet_20240620: true,
		Claude3Opus:             true,
		Claude3Sonnet:           true,
		Claude3Haiku:            true,
		ClaudeV2_1:              true,
	}

	completeCompatibleModels = map[Model]bool{
		Claude35Sonnet:         true,
		Claude3Opus:            true,
		Claude3Sonnet:          true,
		Claude3Haiku:           true,
		ClaudeV2_1:             true,
		ClaudeV2:               true,
		ClaudeV1:               true,
		ClaudeV1_100k:          true,
		ClaudeInstantV1:        true,
		ClaudeInstantV1_100k:   true,
		ClaudeV1_3:             true,
		ClaudeV1_3_100k:        true,
		ClaudeV1_2:             true,
		ClaudeV1_0:             true,
		ClaudeInstantV1_1:      true,
		ClaudeInstantV1_1_100k: true,
		ClaudeInstantV1_0:      true,
	}
)

func (m Model) IsImageCompatible() bool {
	return imageCompatibleModels[m]
}

func (m Model) IsMessageCompatible() bool {
	switch m {
	case Claude3Opus, Claude3Sonnet, Claude3Haiku, ClaudeV2_1, Claude35Sonnet, Claude35Sonnet_20241022, Claude35Sonnet_20240620, Claude35Haiku, Claude35Haiku_20241022:
		return true
	}
	return false
}

func (m Model) IsCompleteCompatible() bool {
	return completeCompatibleModels[m]
}

func (m Model) IsValid() bool {
	_, inImage := imageCompatibleModels[m]
	_, inMessage := messageCompatibleModels[m]
	_, inComplete := completeCompatibleModels[m]
	return inImage || inMessage || inComplete
}
