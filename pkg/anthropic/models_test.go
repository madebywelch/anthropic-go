package anthropic

import (
	"testing"
)

func TestModelCompatibility(t *testing.T) {
	tests := []struct {
		name         string
		model        Model
		wantImage    bool
		wantMessage  bool
		wantComplete bool
		wantValid    bool
	}{
		{"Claude 3.5 Sonnet", Claude35Sonnet, true, true, true, true},
		{"Claude 3 Opus", Claude3Opus, true, true, true, true},
		{"Claude 3 Sonnet", Claude3Sonnet, true, true, true, true},
		{"Claude 3 Haiku", Claude3Haiku, true, true, true, true},
		{"Claude 2.1", ClaudeV2_1, false, true, true, true},
		{"Claude 2", ClaudeV2, false, false, true, true},
		{"Claude 1", ClaudeV1, false, false, true, true},
		{"Claude 1 100k", ClaudeV1_100k, false, false, true, true},
		{"Claude Instant V1", ClaudeInstantV1, false, false, true, true},
		{"Claude Instant V1 100k", ClaudeInstantV1_100k, false, false, true, true},
		{"Claude 1.3", ClaudeV1_3, false, false, true, true},
		{"Claude 1.3 100k", ClaudeV1_3_100k, false, false, true, true},
		{"Claude 1.2", ClaudeV1_2, false, false, true, true},
		{"Claude 1.0", ClaudeV1_0, false, false, true, true},
		{"Claude Instant V1.1", ClaudeInstantV1_1, false, false, true, true},
		{"Claude Instant V1.1 100k", ClaudeInstantV1_1_100k, false, false, true, true},
		{"Claude Instant V1.0", ClaudeInstantV1_0, false, false, true, true},
		{"Invalid Model", Model("NOT A REAL MODEL"), false, false, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.model.IsImageCompatible(); got != tt.wantImage {
				t.Errorf("IsImageCompatible() = %v, want %v", got, tt.wantImage)
			}
			if got := tt.model.IsMessageCompatible(); got != tt.wantMessage {
				t.Errorf("IsMessageCompatible() = %v, want %v", got, tt.wantMessage)
			}
			if got := tt.model.IsCompleteCompatible(); got != tt.wantComplete {
				t.Errorf("IsCompleteCompatible() = %v, want %v", got, tt.wantComplete)
			}
			if got := tt.model.IsValid(); got != tt.wantValid {
				t.Errorf("IsValid() = %v, want %v", got, tt.wantValid)
			}
		})
	}
}

func TestDefaultModel(t *testing.T) {
	if DefaultModel != Claude3Sonnet {
		t.Errorf("DefaultModel = %v, want %v", DefaultModel, Claude3Sonnet)
	}
}
