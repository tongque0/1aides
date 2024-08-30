package planning

import (
	"1aides/pkg/components/db"
	"context"
	"fmt"

	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
)

// 定时任务
// 需要有 类型，时间，内容，发送对象，是否完成，是否删除
// 能够添加定时任务
// 能够删除定时任务

func GoPlanning() {
	// 启用包含秒的cron调度器
	c := cron.New(cron.WithSeconds())

	// 添加一个每秒执行一次的定时任务
	c.AddFunc("* * * * * *", func() {
		fmt.Println("执行定时任务：每秒执行一次")
	})

	// 启动定时任务
	c.Start()

	// 阻塞主线程，直到接收到退出信号
	select {}
}

func InitPlan(c *cron.Cron) {
	collection := db.GetMongoDB().Collection("plantask")
	// 从数据库中获取所有的定时任务
	// 从数据库中获取所有的定时任务
	filter := bson.M{"deleted": false}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		fmt.Println("错误:", err)
		return
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var task PlanTask
		if err := cursor.Decode(&task); err != nil {
			fmt.Println("错误:", err)
			continue
		}
		c.AddFunc(task.TaskTime, func() {
			fmt.Println("执行定时任务：", task.Content)
		})
	}
}

func AddPlan() {
	// 添加定时任务
}

func DeletePlan() {
	// 删除定时任务
}

func UpdatePlan() {
	// 更新定时任务
}
