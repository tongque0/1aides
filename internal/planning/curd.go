package planning

import (
	"1aides/pkg/components/db"
	"1aides/pkg/log/zlog"
	"context"

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
