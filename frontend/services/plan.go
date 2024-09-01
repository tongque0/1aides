package services

import (
	"1aides/internal/planning"

	"github.com/gin-gonic/gin"
)

// HomeHandler 处理主页请求
func PlanHandler(c *gin.Context) {
	c.HTML(200, "plan.tmpl", gin.H{
		"ActivePage": "plan",              // 设置活动页面
		"PlanTasks":  planning.GetPlans(), // 确保这里的数据字段名与模板匹配
	})
}
