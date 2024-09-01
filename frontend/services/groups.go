package services

import (
	"1aides/internal/groups"

	"github.com/gin-gonic/gin"
)

// HomeHandler 处理主页请求
func GroupsHandler(c *gin.Context) {
	c.HTML(200, "groups.tmpl", gin.H{
		"ActivePage": "groups",           // 设置活动页面
		"Groups":     groups.GetGroups(), // 确保这里的数据字段名与模板匹配
	})
}
