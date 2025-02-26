package client

import (
	"context"
	"testing"

	"github.com/pigeonlaser/anthropic-go/v3/pkg/anthropic/client/bedrock"
	"github.com/pigeonlaser/anthropic-go/v3/pkg/anthropic/client/native"
)

func TestMakeClientNativeSuccess(t *testing.T) {
	config := native.Config{
		APIKey: "test",
	}
	ctx := context.Background()
	_, err := MakeClient(ctx, config)
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestMakeClientBedrockSuccess(t *testing.T) {
	config := bedrock.Config{
		Region: "us-west-2",
	}
	ctx := context.Background()
	_, err := MakeClient(ctx, config)
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestMakeClientInvalidConfig(t *testing.T) {
	config := "perhaps-the-archive-is-incomplete"
	ctx := context.Background()
	_, err := MakeClient(ctx, config)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
