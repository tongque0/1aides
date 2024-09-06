package services

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置所有页面和API路由
func SetupRoutes(router *gin.Engine) {
	// 注册登录相关的路由，不需要Auth中间件
	RegisterLoginRoutes(router)

	// 注册受保护的路由组，使用AuthMiddleware中间件
	protected := router.Group("/", AuthMiddleware())

	// 注册各个模块的受保护路由
	RegisterHomeRoutes(protected)
	RegisterPlanRoutes(protected)
	RegisterFriendsRoutes(protected)
	RegisterGroupsRoutes(protected)
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("username")

		if user == nil {
			// 如果Session中没有存储用户信息，重定向到登录页面
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// 用户已登录，继续处理请求
		c.Next()
	}
}
