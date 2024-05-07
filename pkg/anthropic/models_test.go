package anthropic

import "testing"

type modelTest struct {
	model           Model
	imageSupport    bool
	messageSupport  bool
	completeSupport bool
}

func getTestCases() []modelTest {
	return []modelTest{
		{
			model:           Claude3Opus,
			imageSupport:    true,
			messageSupport:  true,
			completeSupport: false,
		},
		{
			model:           Claude3Sonnet,
			imageSupport:    true,
			messageSupport:  true,
			completeSupport: false,
		},
		{
			model:           Claude3Haiku,
			imageSupport:    true,
			messageSupport:  true,
			completeSupport: false,
		},
		{
			model:           ClaudeV2_1,
			imageSupport:    false,
			messageSupport:  true,
			completeSupport: true,
		},
		{
			model:           ClaudeV2,
			imageSupport:    false,
			messageSupport:  false,
			completeSupport: true,
		},
		{
			model:           ClaudeV1,
			imageSupport:    false,
			messageSupport:  false,
			completeSupport: true,
		},
		{
			model:           ClaudeV1_100k,
			imageSupport:    false,
			messageSupport:  false,
			completeSupport: true,
		},
		{
			model:           ClaudeInstantV1,
			imageSupport:    false,
			messageSupport:  false,
			completeSupport: true,
		},
		{
			model:           ClaudeInstantV1_100k,
			imageSupport:    false,
			messageSupport:  false,
			completeSupport: true,
		},
		{
			model:           ClaudeV1_3,
			imageSupport:    false,
			messageSupport:  false,
			completeSupport: true,
		},
		{
			model:           ClaudeV1_3_100k,
			imageSupport:    false,
			messageSupport:  false,
			completeSupport: true,
		},
		{
			model:           ClaudeV1_2,
			imageSupport:    false,
			messageSupport:  false,
			completeSupport: true,
		},
		{
			model:           ClaudeV1_0,
			imageSupport:    false,
			messageSupport:  false,
			completeSupport: true,
		},
		{
			model:           ClaudeInstantV1_1,
			imageSupport:    false,
			messageSupport:  false,
			completeSupport: true,
		},
		{
			model:           ClaudeInstantV1_1_100k,
			imageSupport:    false,
			messageSupport:  false,
			completeSupport: true,
		},
		{
			model:           ClaudeInstantV1_0,
			imageSupport:    false,
			messageSupport:  false,
			completeSupport: true,
		},
		{
			model:           Model("NOT A REAL MODEL"),
			imageSupport:    false,
			messageSupport:  false,
			completeSupport: false,
		},
	}
}

func TestIsImageCompatible(t *testing.T) {
	testCases := getTestCases()
	for _, test := range testCases {
		result := test.model.IsImageCompatible()
		if result != test.imageSupport {
			t.Errorf("IsImageCompatible() for model %s returned %t, expected %t", test.model, result, test.imageSupport)
		}
	}
}
func TestIsMessageCompatible(t *testing.T) {
	testCases := getTestCases()
	for _, test := range testCases {
		result := test.model.IsMessageCompatible()
		if result != test.messageSupport {
			t.Errorf("IsMessageCompatible() for model %s returned %t, expected %t", test.model, result, test.messageSupport)
		}
	}
}
func TestIsCompleteCompatible(t *testing.T) {
	testCases := getTestCases()
	for _, test := range testCases {
		result := test.model.IsCompleteCompatible()
		if result != test.completeSupport {
			t.Errorf("IsCompleteCompatible() for model %s returned %t, expected %t", test.model, result, test.completeSupport)
		}
	}
}
