package anthropic

// Model represents a Claude model.
type Model string

// DefaultModel is the default model used if none is specified.
const DefaultModel = Claude3Sonnet

// https://docs.anthropic.com/claude/docs/models-overview
const (
	// Claude 3 models
	Claude35Sonnet Model = "claude-3-5-sonnet-20240620"
	Claude3Opus    Model = "claude-3-opus-20240229"
	Claude3Sonnet  Model = "claude-3-sonnet-20240229"
	Claude3Haiku   Model = "claude-3-haiku-20240307"

	// Claude 2 models
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
)

var (
	imageCompatibleModels = map[Model]bool{
		Claude35Sonnet: true,
		Claude3Opus:    true,
		Claude3Sonnet:  true,
		Claude3Haiku:   true,
	}

	messageCompatibleModels = map[Model]bool{
		Claude35Sonnet: true,
		Claude3Opus:    true,
		Claude3Sonnet:  true,
		Claude3Haiku:   true,
		ClaudeV2_1:     true,
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

	validModels = make(map[Model]bool)
)

func init() {
	// Populate validModels map
	for model := range imageCompatibleModels {
		validModels[model] = true
	}
	for model := range messageCompatibleModels {
		validModels[model] = true
	}
	for model := range completeCompatibleModels {
		validModels[model] = true
	}
}

func (m Model) IsImageCompatible() bool {
	return imageCompatibleModels[m]
}

func (m Model) IsMessageCompatible() bool {
	return messageCompatibleModels[m]
}

func (m Model) IsCompleteCompatible() bool {
	return completeCompatibleModels[m]
}

func (m Model) IsValid() bool {
	return validModels[m]
}
