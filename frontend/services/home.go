package services

import (
	"1aides/internal/message"
	"sync"

	"github.com/gin-gonic/gin"
)

// 用于存储二维码URL的通道
var qrCodeChan = make(chan string, 1)
var loginOnce sync.Once

// HomeHandler 处理主页请求
func HomeHandler(c *gin.Context) {
	c.HTML(200, "index.tmpl", gin.H{
		"title":    "首页",
		"loginimg": loginimg(),
	})
}

func loginimg() string {
	return message.GetQRCodeURL()
}
