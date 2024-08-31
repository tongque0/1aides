package services

import (
	"1aides/internal/groups"

	"github.com/gin-gonic/gin"
)

// HomeHandler 处理主页请求
func GroupsHandler(c *gin.Context) {
	c.HTML(200, "groups.tmpl", gin.H{
		"ActivePage": "groups",     // 设置活动页面
		"Groups": getGroups(), // 确保这里的数据字段名与模板匹配
	})
}

// getFriends 模拟返回一些好友数据
func getGroups() []groups.Group {
	return []groups.Group{
		{
			ID:            "1",
			HasPermission: true,
			NickName:      "张三",
			RemarkName:    "三哥",
		},
		{
			ID:            "2",
			HasPermission: false,
			NickName:      "李四",
			RemarkName:    "四弟",
		},
	}
}
