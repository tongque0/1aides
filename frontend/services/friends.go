package services

import (
	"1aides/internal/friends"

	"github.com/gin-gonic/gin"
)

// FriendsHandler 处理好友页面请求
func FriendsHandler(c *gin.Context) {
	c.HTML(200, "friends.tmpl", gin.H{
		"ActivePage": "friends",            // 设置活动页面
		"Friends":    friends.GetFriends(), // 确保这里的数据字段名与模板匹配
	})
}
