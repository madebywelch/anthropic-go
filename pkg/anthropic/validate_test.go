package anthropic

import (
	"fmt"
	"testing"
)

type validateMessageTestCase struct {
	request *MessageRequest
	expErr  string
}

func TestValidateMessageRequest(t *testing.T) {
	requests := []validateMessageTestCase{
		{
			request: &MessageRequest{
				Stream: true,
			},
			expErr: "cannot use Message with streaming enabled, use MessageStream instead",
		},
		{
			request: &MessageRequest{
				Stream: false,
				Model:  Model("not-a-valid-model"),
			},
			expErr: "model not-a-valid-model is not compatible with the message endpoint",
		},
		{
			request: &MessageRequest{
				Stream: false,
				Model:  ClaudeV2_1,
				Messages: []MessagePartRequest{{
					Role: "user",
					Content: []ContentBlock{
						NewImageContentBlock(MediaTypeJPEG, "a-gosh-dang-hot-dog"),
					},
				}},
			},
			expErr: fmt.Sprintf("model %s does not support image content", ClaudeV2_1),
		},
		{
			request: &MessageRequest{
				Stream: false,
				Model:  Claude3Opus,
				Messages: []MessagePartRequest{{
					Role:    "user",
					Content: getTwentyOneImgs(),
				}},
			},
			expErr: fmt.Sprintln("too many image content blocks, maximum is 20"),
		},
	}

	for _, test := range requests {
		err := ValidateMessageRequest(test.request)
		if err == nil && test.expErr != "" {
			t.Errorf("Expected error %s, got nil", test.expErr)
		}

		if err == nil {
			continue
		}

		if err.Error() != test.expErr {
			t.Errorf("Expected error %s, got %s", test.expErr, err.Error())
		}
	}
}

func getTwentyOneImgs() []ContentBlock {
	blocks := []ContentBlock{}
	for i := 0; i < 21; i++ {
		blocks = append(
			blocks,
			NewImageContentBlock(MediaTypeJPEG, fmt.Sprintf("a-gosh-dang-hot-dog-%d", i)),
		)
	}

	return blocks
}

func TestValidateMessageStreamRequest(t *testing.T) {
	requests := []validateMessageTestCase{
		{
			request: &MessageRequest{
				Stream: false,
			},
			expErr: "cannot use MessageStream with streaming disabled, use Message instead",
		},
		{
			request: &MessageRequest{
				Stream: true,
				Model:  Model("not-a-valid-model"),
			},
			expErr: "model not-a-valid-model is not compatible with the messagestream endpoint",
		},
		{
			request: &MessageRequest{
				Stream: true,
				Model:  ClaudeV2_1,
				Messages: []MessagePartRequest{{
					Role: "user",
					Content: []ContentBlock{
						NewImageContentBlock(MediaTypeJPEG, "a-gosh-dang-hot-dog"),
					},
				}},
			},
			expErr: fmt.Sprintf("model %s does not support image content", ClaudeV2_1),
		},
		{
			request: &MessageRequest{
				Stream: true,
				Model:  Claude3Opus,
				Messages: []MessagePartRequest{{
					Role:    "user",
					Content: getTwentyOneImgs(),
				}},
			},
			expErr: fmt.Sprintln("too many image content blocks, maximum is 20"),
		},
	}

	for _, test := range requests {
		err := ValidateMessageStreamRequest(test.request)
		if err == nil && test.expErr != "" {
			t.Errorf("Expected error %s, got nil", test.expErr)
		}

		if err == nil {
			continue
		}

		if err.Error() != test.expErr {
			t.Errorf("Expected error %s, got %s", test.expErr, err.Error())
		}
	}
}
