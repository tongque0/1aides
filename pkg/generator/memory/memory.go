package memory

type Memory struct {
	Memory  string              // 记忆内容
	MsgList map[string][]string // 消息列表
}

func (m *Memory) GetMemory() string {
	return m.Memory
}

func (m *Memory) SetMemory(memory string) {
	m.Memory = memory
}

func (m *Memory) AddMsgList(role string, msg string) {
	if m.MsgList == nil {
		m.MsgList = make(map[string][]string)
	}
	m.MsgList[role] = append(m.MsgList[role], msg)
}

func (m *Memory) GetMsgList(role string) []string {
	return m.MsgList[role]
}
