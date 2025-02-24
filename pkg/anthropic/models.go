package anthropic

// Model represents a Claude model.
type Model string

// https://docs.anthropic.com/claude/docs/models-overview
const (
	Claude37Sonnet Model = "claude-3-7-sonnet-latest"

	Claude37Sonnet_20250219 Model = "claude-3-7-sonnet-20250219"

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

	// Most powerful model for highly complex tasks.
	Claude3Opus Model = "claude-3-opus-20240229"

	// Ideal balance of intelligence and speed for enterprise workloads
	Claude3Sonnet Model = "claude-3-sonnet-20240229"

	// Fastest and most compact model for near-instant responsiveness
	Claude3Haiku Model = "claude-3-haiku-20240307"

	// Updated version of Claude 2 with improved accuracy
	ClaudeV2_1 Model = "claude-2.1"

	// Superior performance on tasks that require complex reasoning.
	ClaudeV2 Model = "claude-2"

	// Our largest model, ideal for a wide range of more complex tasks.
	ClaudeV1 Model = "claude-v1"

	// An enhanced version of ClaudeV1 with a 100,000 token context window.
	ClaudeV1_100k Model = "claude-v1-100k"

	// A smaller model with far lower latency, sampling at roughly 40 words/sec!
	ClaudeInstantV1 Model = "claude-instant-v1"

	// An enhanced version of ClaudeInstantV1 with a 100,000 token context window.
	ClaudeInstantV1_100k Model = "claude-instant-v1-100k"

	// Specific sub-versions of the models:

	// More robust against red-team inputs, better at precise instruction-following,
	// better at code, and better and non-English dialogue and writing.
	ClaudeV1_3 Model = "claude-v1.3"

	// An enhanced version of ClaudeV1_3 with a 100,000 token context window.
	ClaudeV1_3_100k Model = "claude-v1.3-100k"

	// An improved version of ClaudeV1, slightly improved at general helpfulness,
	// instruction following, coding, and other tasks. It is also considerably
	// better with non-English languages.
	ClaudeV1_2 Model = "claude-v1.2"

	// An earlier version of ClaudeV1.
	ClaudeV1_0 Model = "claude-v1.0"

	// Latest version of ClaudeInstantV1. It is better at a wide variety of tasks
	// including writing, coding, and instruction following. It performs better on
	// academic benchmarks, including math, reading comprehension, and coding tests.
	ClaudeInstantV1_1 Model = "claude-instant-v1.1"

	// An enhanced version of ClaudeInstantV1_1 with a 100,000 token context window.
	ClaudeInstantV1_1_100k Model = "claude-instant-v1.1-100k"

	// An earlier version of ClaudeInstantV1.
	ClaudeInstantV1_0 Model = "claude-instant-v1.0"
)

func (m Model) IsImageCompatible() bool {
	switch m {
	case Claude3Haiku, Claude3Opus, Claude3Sonnet, Claude35Sonnet, Claude35Sonnet_20241022, Claude35Sonnet_20240620:
		return true
	}
	return false
}

func (m Model) IsMessageCompatible() bool {
	switch m {
	case Claude3Opus, Claude3Sonnet, Claude3Haiku, ClaudeV2_1, Claude35Sonnet, Claude35Sonnet_20241022, Claude35Sonnet_20240620, Claude35Haiku, Claude35Haiku_20241022:
		return true
	}
	return false
}

func (m Model) IsCompleteCompatible() bool {
	switch m {
	case ClaudeV2_1, ClaudeV2, ClaudeV1, ClaudeV1_100k, ClaudeInstantV1, ClaudeInstantV1_100k, ClaudeV1_3, ClaudeV1_3_100k, ClaudeV1_2, ClaudeV1_0, ClaudeInstantV1_1, ClaudeInstantV1_1_100k, ClaudeInstantV1_0:
		return true
	}

	return false
}
