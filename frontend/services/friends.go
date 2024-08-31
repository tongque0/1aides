package services

import (
	"1aides/internal/friends"

	"github.com/gin-gonic/gin"
)

// FriendsHandler 处理好友页面请求
func FriendsHandler(c *gin.Context) {
	c.HTML(200, "friends.tmpl", gin.H{
		"ActivePage": "friends",    // 设置活动页面
		"Friends":    getFriends(), // 确保这里的数据字段名与模板匹配
	})
}

// getFriends 模拟返回一些好友数据
func getFriends() []friends.Friend {
	return []friends.Friend{
		{
			ID:            "1",
			HasPermission: true,
			NickName:      "张三",
			RemarkName:    "三哥",
			Memory:        "老友",
			MsgList: []map[string]string{
				{"role": "user", "content": "你好"},
				{"role": "bot", "content": "你好，有什么可以帮忙的吗？"},
			},
		},
		{
			ID:            "2",
			HasPermission: false,
			NickName:      "李四",
			RemarkName:    "四弟",
			Memory:        "同事",
			MsgList: []map[string]string{
				{"role": "user", "content": "周末有空吗？"},
				{"role": "bot", "content": "这个周末我没空。"},
			},
		},
	}
}
