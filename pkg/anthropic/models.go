package anthropic

// Model represents the model to use for the completion request.
type Model string

const (
	// Our largest model, ideal for a wide range of more complex tasks.
	// Using this model name will automatically switch you to newer versions
	// of claude-v1 after a period of early access evaluation. It currently points at claude-v1.0.
	ClaudeV1 Model = "claude-v1"
	// Current default for claude-v1.
	ClaudeV1_0 Model = "claude-v1.0"
	// [Early access for evaluation] An improved version of claude-v1.
	// It is slightly improved at general helpfulness, instruction following, c
	// oding, and other tasks. It is also considerably better with non-English languages.
	// This model also has the ability to role play (in harmless ways) more consistently,
	// and it defaults to writing somewhat longer and more thorough responses.
	ClaudeV1_2 Model = "claude-v1.2"
	// A smaller model with far lower latency, sampling at roughly 40 words/sec!
	// Its output quality is somewhat lower than claude-v1 models, particularly for
	// complex tasks. However, it is much less expensive and blazing fast. We believe
	// that this model provides more than adequate performance on a range of tasks
	// including text classification, summarization, and lightweight chat applications,
	// as well as search result summarization. Using this model name will automatically
	// switch you to newer versions of claude-instant-v1 after a period of early access evaluation.
	// It currently points at claude-instant-v1.0.
	ClaudeInstantV1 Model = "claude-instant-v1"
	// Current default for claude-instant-v1.
	ClaudeInstantV1_0 Model = "claude-instant-v1.0"
)
