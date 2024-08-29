package bot

import (
	"github.com/eatmoreapple/openwechat"
)

var WxBot *openwechat.Bot

func InitBot() {
	WxBot = openwechat.DefaultBot(openwechat.Desktop)
}
