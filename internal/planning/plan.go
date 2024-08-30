package planning

import (
	"1aides/pkg/components/bot"
	"1aides/pkg/components/db"
	"1aides/pkg/log/zlog"
	"context"
	"fmt"
	"time"

	"github.com/eatmoreapple/openwechat"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// 定时任务
// 需要有 类型，时间，内容，发送对象，是否完成，是否删除
// 能够添加定时任务
// 能够删除定时任务

func GoPlanning() {
	// 启用包含秒的cron调度器
	c := cron.New(cron.WithSeconds())

	// 添加一个每秒执行一次的定时任务
	InitPlanFormDB(c)

	// 启动定时任务
	c.Start()

	// 阻塞主线程，直到接收到退出信号
	select {}
}

func InitPlanFormDB(c *cron.Cron) {
	collection := db.GetMongoDB().Collection("plantask")
	filter := bson.M{"deleted": false}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		zlog.Error("查询失败", zap.Error(err))
		return
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var task PlanTask
		if err := cursor.Decode(&task); err != nil {
			zlog.Error("解码失败", zap.Error(err))
			continue
		}

		var entryID cron.EntryID
		entryID, err = c.AddFunc(task.TaskTime, func() {
			task.EntryID = entryID
			AddPlan(c, collection, task)
		})
		if err != nil {
			zlog.Error("添加定时任务失败", zap.Error(err))
		}

	}
}

// AddPlan 执行任务并根据任务类型标记为完成或更新状态
func AddPlan(c *cron.Cron, collection *mongo.Collection, task PlanTask) {
	// 处理发送消息逻辑
	if task.Recipients != nil {
		for _, recipient := range task.Recipients {
			sendMessage(recipient)
		}
	} 

	// 根据任务类型更新任务状态
	if task.TaskType == SingleTask {
		markTaskAsCompleted(c, collection, task)
	} else if task.TaskType == Recurring {
		updateLastExecuted(collection, task)
	}
}

// sendMessage 发送消息的逻辑
func sendMessage(recipient Object) {
	// 发送消息逻辑
	self, err := bot.WxBot.GetCurrentUser()
	if err != nil {
		zlog.Error("获取当前用户失败", zap.Error(err))
		return
	}
	friends, err := self.Friends()
	if err != nil {
		zlog.Error("获取好友列表失败", zap.Error(err))
		return
	}
	sult := friends.Search(1, func(friend *openwechat.Friend) bool { return friend.ID() == recipient.ID })
	fmt.Println(sult)
	fmt.Printf("发送消息给: %s，类型: %s\n", recipient.ID, recipient.Type)
}

// markTaskAsCompleted 更新任务为已完成并移除定时任务
func markTaskAsCompleted(c *cron.Cron, collection *mongo.Collection, task PlanTask) {
	update := bson.M{
		"$set": bson.M{
			"completed":    true,
			"lastExecuted": time.Now(),
		},
	}

	_, err := collection.UpdateOne(context.TODO(), bson.M{"_id": task.ID}, update)
	if err != nil {
		zlog.Error("更新任务状态失败", zap.Error(err))
	}
	c.Remove(task.EntryID)
}

// updateLastExecuted 仅更新最后执行时间
func updateLastExecuted(collection *mongo.Collection, task PlanTask) {
	update := bson.M{
		"$set": bson.M{
			"lastExecuted": time.Now(),
		},
	}

	_, err := collection.UpdateOne(context.TODO(), bson.M{"_id": task.ID}, update)
	if err != nil {
		zlog.Error("更新任务状态失败", zap.Error(err))
	}
}
