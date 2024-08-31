package services

import (
	"github.com/gin-gonic/gin"
)

// LoginHandler 处理登陆请求
func LoginHandler(c *gin.Context) {
	c.HTML(200, "groups.tmpl", gin.H{
		"ActivePage": "login", // 设置活动页面
	})
}
