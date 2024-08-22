package modhub

import (
	"1aides/pkg/components/generator/memory"
	"1aides/pkg/components/generator/msgchan"
	"1aides/pkg/log/zlog"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

func (m *Model) genGPT(msgchan *msgchan.MsgChan, memory *memory.Memory) {
	config := openai.DefaultConfig(m.Config.APIKey)
	config.BaseURL = m.Config.BaseURL
	c := openai.NewClientWithConfig(config)
	ctx := context.Background()

	req := openai.ChatCompletionRequest{
		Model: m.Config.Model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: m.Config.Prompt,
			},
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: memory.Memory,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: msgchan.Msg.Content,
			},
		},
		Stream: true,
	}
	stream, err := c.CreateChatCompletionStream(ctx, req)
	if err != nil {
		zlog.Warn("流式回复出错", zap.Error(err))
		return
	}
	defer stream.Close()

	for {
		response, err := stream.Recv()

		if errors.Is(err, io.EOF) {
			return
		}

		if err != nil {
			zlog.Warn("流式回复出错", zap.Error(err))
			return
		}
		msgchan.AddMessage(response.Choices[0].Delta.Content)
		fmt.Printf(response.Choices[0].Delta.Content)
	}
}
