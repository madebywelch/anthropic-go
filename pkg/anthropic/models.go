package anthropic

// Model represents a Claude model.
type Model string

const (
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
