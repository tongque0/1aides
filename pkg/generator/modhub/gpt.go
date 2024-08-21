package modhub

import (
	"1aides/pkg/generator/memory"
	"1aides/pkg/generator/msgchan"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/sashabaranov/go-openai"
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
		fmt.Printf("GPT生成错误: %v\n", err)
		return
	}
	defer stream.Close()

	for {
		response, err := stream.Recv()

		if errors.Is(err, io.EOF) {
			return
		}

		if err != nil {
			fmt.Printf("\nStream error: %v\n", err)
			return
		}
		msgchan.AddMessage(response.Choices[0].Delta.Content)
		fmt.Printf(response.Choices[0].Delta.Content)
	}
}
