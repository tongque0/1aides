package message

import (
	"1aides/internal/helper"
	"1aides/pkg/components/auth"

	"github.com/eatmoreapple/openwechat"
)

func HandleMessage(msg *openwechat.Message) {

	if helper.IsHelp(msg) {
		go helper.Helper(msg)
		return
	}

	if auth.Auth(msg) {
		if msg.IsText() {
			go gen(msg)
		}
	}
}

func HandleUUID(uuid string) string {
	qrcodeUrl := openwechat.GetQrcodeUrl(uuid)
	return qrcodeUrl
}
