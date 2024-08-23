package message

import (
	"1aides/internal/helper"
	"1aides/pkg/components/auth"

	"github.com/eatmoreapple/openwechat"
)

func HandleMessage(msg *openwechat.Message) {

	if helper.IsHelp(msg) {
		helper.Helper(msg)
		return
	}

	if auth.Auth(msg) {
		if msg.IsText() {
			gen(msg)
		}
	}
}
