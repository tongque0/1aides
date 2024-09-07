package planning

import (
	"1aides/pkg/components/bot"
	"1aides/pkg/components/db"
	"1aides/pkg/log/zlog"
	"context"
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
var C *cron.Cron

func GoPlanning() {
	// 启用包含秒的cron调度器
	C = cron.New(cron.WithSeconds())

	// 添加一个每秒执行一次的定时任务
	initPlanFormDB(C)

	// 启动定时任务
	C.Start()

	// 阻塞主线程，直到接收到退出信号
	select {}
}

func initPlanFormDB(c *cron.Cron) {
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
			addPlan(c, collection, task)
		})
		if err != nil {
			zlog.Error("添加定时任务失败", zap.Error(err))
		}

	}
}

// AddPlan 执行任务并根据任务类型标记为完成或更新状态
func addPlan(c *cron.Cron, collection *mongo.Collection, task PlanTask) {
	// 处理发送消息逻辑
	if task.Recipients != nil {
		sendMessage(task.Recipients, task.Content)
	}

	// 根据任务类型更新任务状态
	if task.TaskType == SingleTask {
		markTaskAsCompleted(c, collection, task)
	} else if task.TaskType == Recurring {
		updateLastExecuted(collection, task)
	}
}

// sendMessage 发送消息的逻辑
func sendMessage(recipient []Object, content Content) {
	// 发送消息逻辑
	self, err := bot.WxBot.GetCurrentUser()
	if err != nil {
		zlog.Error("定时任务获取机器人失败", zap.Error(err))
		return
	}
	friends, err := self.Friends()
	if err != nil {
		zlog.Error("定时任务获取好友列表失败", zap.Error(err))
		return
	}
	groups, err := self.Groups()
	if err != nil {
		zlog.Error("定时任务群组列表失败", zap.Error(err))
		return
	}
	for _, r := range recipient {
		if r.Type == User {
			friend := friends.Search(1, func(friend *openwechat.Friend) bool { return friend.ID() == r.ID })[0]
			friend.SendText(content.Detail)
		} else if r.Type == Group {
			group := groups.Search(1, func(group *openwechat.Group) bool { return group.ID() == r.ID })[0]
			group.SendText(content.Detail)
		} else if r.Type == FileHelper {
			self.FileHelper().SendText(content.Detail)
		}
	}
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
