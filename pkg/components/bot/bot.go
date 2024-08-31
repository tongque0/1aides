package bot

import (
	"1aides/pkg/log/zlog"

	"github.com/eatmoreapple/openwechat"
	"go.uber.org/zap"
)

var WxBot *openwechat.Bot
var globalQRCodeURL string

func InitBot() {
	WxBot = nil
	WxBot = openwechat.DefaultBot(openwechat.Desktop)
	WxBot.UUIDCallback = handleUUID
	WxBot.ScanCallBack = handleScan
	WxBot.LoginCallBack = handLogin
}

func handleUUID(uuid string) {
	globalQRCodeURL = openwechat.GetQrcodeUrl(uuid)
}

func GetLoginURL() string {
	return globalQRCodeURL
}

func handleScan(user openwechat.CheckLoginResponse) {
	avatar, err := user.Avatar()
	if err != nil {
		globalQRCodeURL = ""
		zlog.Error("Failed to get avatar", zap.Error(err))
		return
	}

	globalQRCodeURL = avatar
}

func handLogin(user openwechat.CheckLoginResponse) {
	avatar, err := user.Avatar()
	if err != nil {
		globalQRCodeURL = ""
		zlog.Error("Failed to get avatar", zap.Error(err))
		return
	}
	globalQRCodeURL = avatar
}
