package services

import (
	"1aides/internal/groups"
	"1aides/pkg/log/zlog"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RegisterGroupsRoutes(router *gin.RouterGroup) {
	// 好友管理页面
	router.GET("/groups", GroupsHandler)
	router.POST("/groups/setpermission", SetGroupPermissionHandler)
}

// HomeHandler 处理主页请求
func GroupsHandler(c *gin.Context) {
	c.HTML(200, "groups.tmpl", gin.H{
		"ActivePage": "groups",           // 设置活动页面
		"Groups":     groups.GetGroups(), // 确保这里的数据字段名与模板匹配
	})
}

// SetPermissionHandler 设置好友权限
func SetGroupPermissionHandler(c *gin.Context) {
	var Permission struct {
		ID         string `json:"id"`
		Permission string `json:"permission"`
	}
	// 从请求体中解析 JSON 数据
	if err := c.BindJSON(&Permission); err != nil {
		zlog.Error("无法解析登录请求的JSON数据", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	id := Permission.ID

	permission := Permission.Permission

	// 解析权限值，支持两种权限：普通权限和管理员权限
	hasPermission := false

	if permission == "admin" {
		hasPermission = true
	} else if permission == "normal" {
		hasPermission = true
	}

	// 调用 friends 包的 SetPermission 函数更新权限
	err := groups.SetPermission(id, hasPermission)

	if err != nil {
		zlog.Error("无法设置好友权限", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(200, gin.H{"status": "success"})
}
