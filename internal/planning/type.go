package planning

import (
	"time"
)

// 对象类型
type ObjectType string
type TaskType string

const (
	SingleTask TaskType = "single"    // 单次任务
	Recurring  TaskType = "recurring" // 循环任务
)
const (
	User       ObjectType = "user"
	Group      ObjectType = "group"
	FileHelper ObjectType = "filehelper"
)

// 对象
type Object struct {
	Type ObjectType // 对象类型，如用户或群组
	ID   string     // 对象ID，例如用户ID、群组ID等
}

// 定时任务结构体
type PlanTask struct {
	TaskType     TaskType // 任务类型
	TaskTime     string
	Content      string    // 任务内容
	Recipients   []Object  // 接收任务通知的对象列表
	LastExecuted time.Time // 最后一次执行时间
	Deleted      bool      // 是否已删除
}
