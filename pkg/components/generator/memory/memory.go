package memory

type Memory struct {
	Memory  string // 记忆内容
	MsgList []map[string]string
}

func (m *Memory) GetMemory() string {
	return m.Memory
}

func (m *Memory) SetMemory(memory string) {
	m.Memory = memory
}

func (m *Memory) AddMsgList(role string, msg string) {
	message := map[string]string{
		"role":    role,
		"content": msg,
	}
	m.MsgList = append(m.MsgList, message)
}

func (m *Memory) GetMsgList() []map[string]string {
	return m.MsgList
}
