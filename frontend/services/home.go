package services

import (
	"1aides/pkg/components/bot"
	"html/template"

	"github.com/gin-gonic/gin"
)

func RegisterHomeRoutes(router *gin.RouterGroup) {
	// 首页路由
	router.GET("/", HomeHandler)
}

// HomeHandler 处理主页请求
func HomeHandler(c *gin.Context) {
	c.HTML(200, "home.tmpl", gin.H{
		"ActivePage": "home", // 设置活动页面
		"loginimg":   template.URL(loginimg()),
	})
}

func loginimg() string {
	return bot.GetLoginURL()
}
