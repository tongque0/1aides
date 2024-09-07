package modhub

import (
	"1aides/pkg/components/generator/memory"
	"1aides/pkg/components/generator/msgchan"
)

type ModelType string

const (
	GPT  ModelType = "GPT"
	BERT ModelType = "BERT"
	T5   ModelType = "T5"
)

type ModelConfig struct {
	Model   string `yaml:"model"`
	APIKey  string `yaml:"apiKey"`
	BaseURL string `yaml:"baseURL"`
	Prompt  string `yaml:"prompt"`
}

type Model struct {
	Type   ModelType   `yaml:"type"`
	Config ModelConfig `yaml:"config"`
}

func NewModel(modelType ModelType, config *ModelConfig) Model {
	if config == nil {
		config = &ModelConfig{
			APIKey:  "sk-29sMaKDD5aBgDtyx02014694972846Cc8c8b9fEb18192532",
			BaseURL: "https://prime.zetatechs.com/v1",
			Prompt:  "你的身份是一位微信消息机器人，你的开发者是同阙。你可以回复任何你想回复的内容，但是要有逻辑。",
		}
	}
	return Model{
		Type:   modelType,
		Config: *config,
	}
}
func NewModelWithString(modelType ModelType, model string, apiKey string, baseURL string, prompt string) Model {
	return Model{
		Type: modelType,
		Config: ModelConfig{
			Model:   model,
			APIKey:  apiKey,
			BaseURL: baseURL,
			Prompt:  prompt,
		},
	}
}
func (m *Model) Gen(msgchan *msgchan.MsgChan, mermory *memory.Memory) {
	switch m.Type {
	case GPT:
		m.genGPT(msgchan, mermory)
	case BERT:
		// do something
	case T5:
		// do something
	default:
		m.genGPT(msgchan, mermory)
	}
}

func (m *Model) GenMemory(msgchan *msgchan.MsgChan, memory *memory.Memory) {
	switch m.Type {
	case GPT:
		m.genMemoryForGPT(msgchan, memory)
	case BERT:
		// do something
	case T5:
		// do something
	}
}
