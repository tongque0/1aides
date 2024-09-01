package services

import (
	"1aides/pkg/components/bot"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置所有页面和API路由
func SetupRoutes(router *gin.Engine) {

	// 对需要认证的路由使用中间件
	protected := router.Group("/", AuthMiddleware())

	// 首页（不受AuthMiddleware影响）
	router.GET("/", HomeHandler)

	// 登陆（不受AuthMiddleware影响）
	router.GET("/login", LoginHandler)

	// 受保护的路由组
	protected.GET("/plan", PlanHandler)
	protected.GET("/friends", FriendsHandler)
	protected.GET("/groups", GroupsHandler)
}

// AuthMiddleware 是一个中间件，用于检查用户是否已登录并且bot.alive为真
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查 bot.alive 的状态
		if !bot.WxBot.Alive() {
			// 如果 bot.alive 为假，跳转到认证页
			c.Redirect(http.StatusTemporaryRedirect, "/")
			c.Abort()
			return
		}

		// 如果验证通过，允许请求继续
		c.Next()
	}
}
