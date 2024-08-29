package bot

import (
	"github.com/eatmoreapple/openwechat"
)

var WxBot *openwechat.Bot
var globalQRCodeURL string

func InitBot() {
	WxBot = openwechat.DefaultBot(openwechat.Desktop)
	WxBot.UUIDCallback = handleUUID
}

func handleUUID(uuid string) {
	globalQRCodeURL = openwechat.GetQrcodeUrl(uuid)
}

func GetLoginURL() string {
	return globalQRCodeURL
}
