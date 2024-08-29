package message

import (
	"1aides/internal/helper"
	"1aides/pkg/components/auth"
	"fmt"

	"github.com/eatmoreapple/openwechat"
)

var globalQRCodeURL string

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

func HandleUUID(uuid string) {
	globalQRCodeURL = openwechat.GetQrcodeUrl(uuid)
	fmt.Println(globalQRCodeURL)
}

func GetQRCodeURL() string {
	return globalQRCodeURL
}
