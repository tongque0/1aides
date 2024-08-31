package services

import (
	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置所有页面和API路由
func SetupRoutes(router *gin.Engine) {
	// 首页
	router.GET("/", HomeHandler)
	// 登陆
	router.GET("/login", LoginHandler)

}
