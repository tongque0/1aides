package modhub

import (
	"1aides/pkg/components/generator/memory"
	"1aides/pkg/components/generator/msgchan"
	"1aides/pkg/log/zlog"
	"context"
	"errors"
	"io"

	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

func (m *Model) genGPT(msgchan *msgchan.MsgChan, memory *memory.Memory) {
	config := openai.DefaultConfig(m.Config.APIKey)
	config.BaseURL = m.Config.BaseURL
	c := openai.NewClientWithConfig(config)
	ctx := context.Background()

	// 初始化请求消息列表
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: m.Config.Prompt,
		},
	}

	// 将 memory.Memory 拼接进去
	if memory.Memory != "" {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: memory.Memory,
		})
	}

	// 将 memory.MsgList 中的每个消息拼接进去
	for _, msg := range memory.GetMsgList() {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    msg["role"],
			Content: msg["content"],
		})
	}

	// 将用户的当前消息拼接进去
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: msgchan.Msg.Content,
	})

	req := openai.ChatCompletionRequest{
		Model:    m.Config.Model,
		Messages: messages,
		Stream:   true,
	}

	stream, err := c.CreateChatCompletionStream(ctx, req)
	if err != nil {
		msgchan.AddMessage("回复发生未知错误")
		zlog.Warn("流式回复出错", zap.Error(err))
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
	}
}

func (m *Model) genMemoryForGPT(msgchan *msgchan.MsgChan, memory *memory.Memory) {
	config := openai.DefaultConfig(m.Config.APIKey)
	config.BaseURL = m.Config.BaseURL
	c := openai.NewClientWithConfig(config)
	ctx := context.Background()

	// 初始化请求消息列表
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "你是一个微信机器聊天机器人记忆模块，你负责记忆的生成和更新。你需要使用下面的现有记忆以及用户发给你的各种聊天内容，更新现有的记忆，记忆需要简洁明了，不要包含无关信息。下面是现有记忆内容：",
		},
	}

	// 将 memory.Memory 拼接进去
	if memory.Memory != "" {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: memory.Memory,
		})
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: "下面是用户发给你的聊天内容：",
		})
	}

	// 将 memory.MsgList 中的每个消息拼接进去
	for _, msg := range memory.GetMsgList() {
		if msg["role"] == "user" {
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    msg["role"],
				Content: msg["content"],
			})
		}
	}

	resp, err := c.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    m.Config.Model,
			Messages: messages,
		},
	)
	if err != nil {
		zlog.Warn("生成记忆出错", zap.Error(err))
		return
	}
	msgchan.AddMemory(resp.Choices[0].Message.Content)
}
