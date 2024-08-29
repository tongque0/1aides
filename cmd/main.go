package main

import (
	"1aides/frontend/services"
	"1aides/internal/message"
	"1aides/pkg/components/bot"
	"1aides/pkg/log/zlog"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {

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

	// 启动登录过程
	go ensureLoggedIn()

	router.Run(":8999")
}

// ensureLoggedIn 确保始终有账号登录
func ensureLoggedIn() {
	for {
		bot.InitBot()
		bot.WxBot.MessageHandler = message.HandleMessage
		bot.WxBot.UUIDCallback = message.HandleUUID
		if err := bot.WxBot.Login(); err != nil {
			zlog.Error("登陆失败，正在重试...", zap.Error(err))
			continue
		}
		zlog.Info("登陆成功")
		bot.WxBot.Block()
	}
}
