package msgchan

import (
	"strings"

	"github.com/eatmoreapple/openwechat"
)

type MsgChan struct {
	builder   strings.Builder
	limitFunc func(*strings.Builder) bool
	Msg       *openwechat.Message
}

func NewMsgChan(limitFunc func(*strings.Builder) bool) *MsgChan {
	if limitFunc == nil {
		limitFunc = defaultLimitFunc
	}
	return &MsgChan{
		builder:   strings.Builder{},
		limitFunc: limitFunc,
		Msg:       nil,
	}
}

func defaultLimitFunc(b *strings.Builder) bool {
	return b.Len() > 500
}

func (m *MsgChan) AddMessage(msg string) {
	m.builder.WriteString(msg)
	if m.limitFunc(&m.builder) {
		m.sendTextData()
	}
}

func (m *MsgChan) Flush() {
	m.sendTextData()
}

func (m *MsgChan) sendTextData() {
	if m.Msg != nil && m.builder.Len() > 0 {
		m.Msg.ReplyText(m.builder.String())
		m.builder.Reset()
	}
}

func (m *MsgChan) Show() string {
	return m.builder.String()
}

func (m *MsgChan) Consume(msg *openwechat.Message) {
	m.Msg = msg
}
