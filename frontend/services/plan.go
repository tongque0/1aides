package services

import (
	"1aides/internal/planning"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// HomeHandler 处理主页请求
func PlanHandler(c *gin.Context) {
	c.HTML(200, "plan.tmpl", gin.H{
		"ActivePage": "plan",     // 设置活动页面
		"PlanTasks":  getPlans(), // 确保这里的数据字段名与模板匹配
	})
}

func getPlans() []planning.PlanTask {
	return []planning.PlanTask{
		{
			ID:           primitive.NewObjectID(),
			TaskType:     planning.SingleTask,
			TaskTime:     "0 0 1 * *", // 每月1号午夜
			Content:      planning.Content{Type: string(planning.Text), Detail: "每月报告"},
			Recipients:   []planning.Object{{Type: planning.User, ID: "user1"}},
			LastExecuted: time.Now().Add(-30 * 24 * time.Hour), // 30天前
			Completed:    false,
			Deleted:      false,
		},
		{
			ID:           primitive.NewObjectID(),
			TaskType:     planning.Recurring,
			TaskTime:     "0 0 * * 1", // 每周一午夜
			Content:      planning.Content{Type: string(planning.Image), Detail: "每周图像"},
			Recipients:   []planning.Object{{Type: planning.Group, ID: "group1"}},
			LastExecuted: time.Now().Add(-7 * 24 * time.Hour), // 7天前
			Completed:    false,
			Deleted:      false,
		},
		{
			ID:           primitive.NewObjectID(),
			TaskType:     planning.SingleTask,
			TaskTime:     "0 12 25 12 *", // 每年12月25日中午
			Content:      planning.Content{Type: string(planning.File), Detail: "年度文件"},
			Recipients:   []planning.Object{{Type: planning.FileHelper, ID: "filehelper1"}},
			LastExecuted: time.Now().Add(-365 * 24 * time.Hour), // 365天前
			Completed:    false,
			Deleted:      false,
		},
	}
}
