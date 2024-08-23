package main

import (
	"1aides/internal/friends"
	"1aides/internal/groups"
	"1aides/internal/message"
	"1aides/pkg/components/bot"
	"1aides/pkg/log/zlog"

	"github.com/eatmoreapple/openwechat"
	"go.uber.org/zap"
)

func main() {
	bot := bot.WxBot

	// 注册消息处理函数
	bot.MessageHandler = message.HandleMessage
	// 注册登陆二维码回调
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl

	// 登陆
	if err := bot.Login(); err != nil {
		zlog.Error("登陆失败", zap.Error(err))
		return
	}
	friends.InitFriendDB()
	groups.InitGroupsDB()

	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	bot.Block()
}
