package message

import (
	"1aides/internal/helper"

	"github.com/eatmoreapple/openwechat"
)

func HandleMessage(msg *openwechat.Message) {

	if helper.IsHelp(msg) {
		go helper.Helper(msg)
		return
	}

	if helper.Auth(msg) {
		if msg.IsText() {
			go gen(msg)
		}
	}
}
