package services

import (
	"1aides/internal/planning"

	"github.com/gin-gonic/gin"
)

func RegisterPlanRoutes(router *gin.RouterGroup) {
	// 计划页面
	router.GET("/plan", PlanHandler)
	// 添加计划
	router.POST("/plan/add", AddPlanHandler)
	// 删除计划
	// router.POST("/plan/delete", DeletePlanHandler)
}

// HomeHandler 处理主页请求
func PlanHandler(c *gin.Context) {
	c.HTML(200, "plan.tmpl", gin.H{
		"ActivePage": "plan",              // 设置活动页面
		"PlanTasks":  planning.GetPlans(), // 确保这里的数据字段名与模板匹配
	})
}

func AddPlanHandler(c *gin.Context) {
	// 构造计划任务数据
	// taskType := c.PostForm("taskType")
	taskTime := c.PostForm("taskTime")
	contentType := c.PostForm("contentType")
	contentDetail := c.PostForm("contentDetail")
	// recipients := c.PostForm("recipients") // 假设以逗号分隔的字符串

	// 创建任务对象
	task := planning.PlanTask{
		TaskType: "single", // 默认为单次任务
		TaskTime: taskTime,
		Content: planning.Content{
			Type:   contentType,
			Detail: contentDetail,
		},
		// Recipients: strings.Split(recipients, ","), // 分割字符串为数组
	}

	// 添加任务到系统
	err := planning.AddPlan(task)
	if err != nil {
		c.JSON(400, gin.H{"error": "无法添加计划任务"})
		return
	}

	// 重定向到计划任务页面
	c.Redirect(302, "/plan")
}

// func DeletePlanHandler(c *gin.Context) {
// 	// 获取表单数据
// 	task := c.PostForm("task")
// 	// 添加任务
// 	planning.AddPlan(task)
// 	// 重定向到计划页面
// 	c.Redirect(302, "/plan")
// }
