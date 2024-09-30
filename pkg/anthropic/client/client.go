package client

import (
	"context"
	"fmt"

	"github.com/madebywelch/anthropic-go/v4/pkg/anthropic"

	"github.com/madebywelch/anthropic-go/v4/pkg/anthropic/client/bedrock"
	"github.com/madebywelch/anthropic-go/v4/pkg/anthropic/client/native"
)

type ClientType string

type Client interface {
	Message(context.Context, *anthropic.MessageRequest) (*anthropic.MessageResponse, error)
	MessageStream(context.Context, *anthropic.MessageRequest) (<-chan *anthropic.MessageStreamResponse, <-chan error)
	Complete(context.Context, *anthropic.CompletionRequest) (*anthropic.CompletionResponse, error)
	CompleteStream(context.Context, *anthropic.CompletionRequest) (<-chan *anthropic.StreamResponse, <-chan error)
}

func MakeClient(ctx context.Context, config interface{}) (Client, error) {
	switch cfg := config.(type) {
	case bedrock.Config:
		return bedrock.MakeClient(ctx, cfg)
	case native.Config:
		return native.MakeClient(cfg)
	}

	return nil, fmt.Errorf("unknown client config")
}
