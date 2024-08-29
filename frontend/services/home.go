package services

import (
	"github.com/gin-gonic/gin"
)

// HomeHandler 处理主页请求
func HomeHandler(c *gin.Context) {
	c.HTML(200, "index.tmpl", gin.H{"title": "首页"})
}
