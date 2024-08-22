package bot

import "github.com/eatmoreapple/openwechat"

var WxBot *openwechat.Bot

func init() {
	WxBot = openwechat.DefaultBot(openwechat.Desktop)
}
