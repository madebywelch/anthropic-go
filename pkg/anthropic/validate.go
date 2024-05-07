package anthropic

import "fmt"

func ValidateMessageRequest(req *MessageRequest) error {
	if req.Stream {
		return fmt.Errorf("cannot use Message with streaming enabled, use MessageStream instead")
	}

	if !req.Model.IsMessageCompatible() {
		return fmt.Errorf("model %s is not compatible with the message endpoint", req.Model)
	}

	if !req.Model.IsImageCompatible() && req.ContainsImageContent() {
		return fmt.Errorf("model %s does not support image content", req.Model)
	}

	if req.CountImageContent() > 20 {
		return fmt.Errorf("too many image content blocks, maximum is 20")
	}

	return nil
}

func ValidateMessageStreamRequest(req *MessageRequest) error {
	if !req.Stream {
		return fmt.Errorf("cannot use MessageStream with streaming disabled, use Message instead")
	}

	if !req.Model.IsMessageCompatible() {
		return fmt.Errorf("model %s is not compatible with the messagestream endpoint", req.Model)
	}

	if !req.Model.IsImageCompatible() && req.ContainsImageContent() {
		return fmt.Errorf("model %s does not support image content", req.Model)
	}

	if req.CountImageContent() > 20 {
		return fmt.Errorf("too many image content blocks, maximum is 20")
	}

	return nil
}

func ValidateCompleteRequest(req *CompletionRequest) error {
	if req.Stream {
		return fmt.Errorf("cannot use Complete with streaming enabled, use CompleteStream instead")
	}

	if !req.Model.IsCompleteCompatible() {
		return fmt.Errorf("model %s is not compatible with the completion endpoint", req.Model)
	}

	return nil
}

func ValidateCompleteStreamRequest(req *CompletionRequest) error {
	if !req.Stream {
		return fmt.Errorf("cannot use CompleteStream with streaming disabled, use Complete instead")
	}

	if !req.Model.IsCompleteCompatible() {
		return fmt.Errorf("model %s is not compatible with the completion endpoint", req.Model)
	}

	return nil
}
