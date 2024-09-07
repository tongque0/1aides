package planning

import (
	"1aides/pkg/components/db"
	"1aides/pkg/log/zlog"
	"context"

	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

func GetPlans() []PlanTask {

	collection := db.GetMongoDB().Collection("plantask")
	filter := bson.M{"deleted": false}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		zlog.Error("查询失败", zap.Error(err))
		return nil
	}
	defer cursor.Close(context.TODO())

	var tasks []PlanTask
	for cursor.Next(context.TODO()) {
		var task PlanTask
		if err := cursor.Decode(&task); err != nil {
			zlog.Error("解码失败", zap.Error(err))
			continue
		}
		tasks = append(tasks, task)
	}
	return tasks
}

func AddPlan(task PlanTask) error {
	collection := db.GetMongoDB().Collection("plantask")
	var entryID cron.EntryID
	entryID, err := C.AddFunc(task.TaskTime, func() {
		task.EntryID = entryID
		addPlan(C, collection, task)
	})

	task.EntryID = entryID
	if err != nil {
		zlog.Error("添加定时任务失败", zap.Error(err))
		return err
	}

	_, err = collection.InsertOne(context.TODO(), task)
	if err != nil {
		zlog.Error("插入失败", zap.Error(err))
		return err
	}
	return nil
}

func DeletedPlan(id string) error {
	collection := db.GetMongoDB().Collection("plantask")
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"deleted": true}}
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		zlog.Error("删除失败", zap.Error(err))
		return err
	}
	return nil
}
