package message

import (
	"github.com/eatmoreapple/openwechat"
)

func HandleMessage(msg *openwechat.Message) {
	if msg.IsText() {
		gen(msg)
	}
}
