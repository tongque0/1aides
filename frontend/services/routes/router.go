package routes

import (
	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置所有页面和API路由
func SetupRoutes(router *gin.Engine) {
	// 首页
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{"title": "首页"})
	})

	// 健康检查接口
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// 其他路由可以在这里继续添加
}
