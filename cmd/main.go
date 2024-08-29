package main

import (
	"1aides/frontend/services"
	"1aides/internal/friends"
	"1aides/internal/groups"
	"1aides/internal/message"
	"1aides/pkg/components/bot"
	"1aides/pkg/log/zlog"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	bot.WxBot.MessageHandler = message.HandleMessage
	// 注册登陆二维码回调
	bot.WxBot.UUIDCallback = message.HandleUUID

	router := gin.Default()

	templatePath := os.Getenv("TEMPLATE_PATH")
	staticPath := os.Getenv("STATIC_PATH")

	if templatePath == "" {
		templatePath = "../frontend/templates/*" // 默认值，适用于开发环境
	}
	if staticPath == "" {
		staticPath = "../frontend/static" // 默认值，适用于开发环境
	}

	router.LoadHTMLGlob(templatePath)

	// Serve static files
	router.Static("/static", staticPath)
	services.SetupRoutes(router)

	// 登陆
	go func() {
		if err := bot.WxBot.Login(); err != nil {
			zlog.Error("登陆失败", zap.Error(err))
			return
		}
		friends.InitFriendDB()
		groups.InitGroupsDB()
	}()

	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	// bot.WxBot.Block()
	router.Run(":8999")
}
