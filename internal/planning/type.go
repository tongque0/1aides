package planning

import (
	"time"

	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 对象类型
type ObjectType string
type TaskType string
type MessageType string

const (
	SingleTask TaskType = "single"    // 单次任务
	Recurring  TaskType = "recurring" // 循环任务
)
const (
	User       ObjectType = "user"
	Group      ObjectType = "group"
	FileHelper ObjectType = "filehelper"
)
const (
	Text  MessageType = "text"
	Image MessageType = "image"
	File  MessageType = "file"
	Voice MessageType = "voice"
)

// 对象
type Object struct {
	Type ObjectType // 对象类型，如用户或群组
	ID   string     // 对象ID，例如用户ID、群组ID等
}

// Content
type Content struct {
	Type   string `bson:"type" json:"type"`     // 内容类型
	Detail string `bson:"detail" json:"detail"` // 内容详细信息
}

// 定时任务结构体
// 定时任务结构体
type PlanTask struct {
	ID           primitive.ObjectID `bson:"_id"`          // 任务ID，MongoDB中的文档ID
	TaskType     TaskType           `bson:"tasktype"`     // 任务类型
	TaskTime     string             `bson:"tasktime"`     // 定时任务时间（Cron 表达式或其他表示方式）
	Content      Content            `bson:"content"`      // 任务内容
	Recipients   []Object           `bson:"recipients"`   // 接收任务通知的对象列表
	LastExecuted time.Time          `bson:"lastExecuted"` // 最后一次执行时间
	Completed    bool               `bson:"completed"`    // 是否已完成
	Deleted      bool               `bson:"deleted"`      // 是否已删除
	EntryID      cron.EntryID       `bson:"entryID"`      // cron.EntryID
}
