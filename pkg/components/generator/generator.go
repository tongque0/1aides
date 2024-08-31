package generator

import (
	"1aides/pkg/components/generator/memory"
	"1aides/pkg/components/generator/modhub"
	"1aides/pkg/components/generator/msgchan"

	"github.com/eatmoreapple/openwechat"
)

type Generator struct {
	Model   modhub.Model     // 消息生成器使用模型
	MsgChan *msgchan.MsgChan // 消息生成器使用的消息通道
	Memory  memory.Memory    // 消息生成器的上下文信息
}

func NewGenerator(msg *openwechat.Message, options ...func(*Generator)) *Generator {
	gen := &Generator{
		Model:   DefaultModel(),
		Memory:  DefaultMemory(),
		MsgChan: DefaultMsgChan(),
	}

	for _, option := range options {
		option(gen)
	}
	gen.MsgChan.Consume(msg)
	return gen
}

// Generate 生成消息
func (g *Generator) Generate() string {
	g.Model.Gen(g.MsgChan, &g.Memory)

	g.MsgChan.Flush()

	return g.MsgChan.GetRecords()
}

// Generate 生成消息
func (g *Generator) GenMemory() string {
	g.Model.GenMemory(g.MsgChan, &g.Memory)

	return g.MsgChan.GetRecords()
}

// DefaultModel 默认模型
func DefaultModel() modhub.Model {
	return modhub.NewModel(modhub.GPT, nil) // Use default configuration
}

// DefaultMemory 默认上下文
func DefaultMemory() memory.Memory {
	return memory.Memory{
		Memory:  "这是你的记忆，作为一个机器人，你的开发者，为你赋予了记忆部分，此记忆会随着你跟我们的相处而不断地增长,该记忆会影响你的回答，所以请你好好的记住这些信息。",
		MsgList: make([]map[string]string, 0),
	}
}

// DefaultMsgChan 默认消息通道
func DefaultMsgChan() *msgchan.MsgChan {
	return msgchan.NewMsgChan(nil)
}

// WithModel 设置模型
func WithModel(model modhub.Model) func(*Generator) {
	return func(g *Generator) { g.Model = model }
}

// WithContext 设置上下文
func WithMemory(memory memory.Memory) func(*Generator) {
	return func(g *Generator) { g.Memory = memory }
}

// WithMsgChan 设置消息通道
func WithMsgChan(msgchan *msgchan.MsgChan) func(*Generator) {
	return func(g *Generator) { g.MsgChan = msgchan }
}
