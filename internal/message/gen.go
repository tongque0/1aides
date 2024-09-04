package message

import (
	"1aides/internal/friends"
	"1aides/pkg/components/db"
	"1aides/pkg/components/generator"
	"1aides/pkg/components/generator/memory"
	"1aides/pkg/components/generator/modhub"
	"1aides/pkg/log/zlog"
	"context"

	"github.com/eatmoreapple/openwechat"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func gen(msg *openwechat.Message) {
	sender, err := msg.Sender()
	if err != nil {
		return
	}

	// 从配置中心加载模型配置
	var model modhub.Model
	modelctl := db.GetMongoDB().Collection("models")
	filter := bson.M{"type": "gpt"}
	err = modelctl.FindOne(context.TODO(), filter).Decode(&model)
	if err != nil {
		zlog.Error("加载配置失败", zap.Error(err))
		msg.ReplyText("模型配置加载失败")
		return
	}
	// 从mongoDB中加载记忆内容
	var memory memory.Memory
	var memoryctl *mongo.Collection
	if sender.IsGroup() {
		memoryctl = db.GetMongoDB().Collection("groups")
	} else {
		memoryctl = db.GetMongoDB().Collection("friends")
	}
	filter = bson.M{"id": sender.ID()}
	var friend friends.Friend
	err = memoryctl.FindOne(context.TODO(), filter).Decode(&friend)
	if err != nil {
		zlog.Warn("未找到对应的好友记录", zap.String("ID", sender.ID()), zap.Error(err))
		friend = friends.Friend{
			Memory:  "",
			MsgList: nil,
		}
	}
	memory.Memory = friend.Memory
	memory.MsgList = friend.MsgList

	// 使用加载的模型配置创建生成器实例
	gen := generator.NewGenerator(msg, generator.WithModel(model), generator.WithMemory(memory))

	result := gen.Generate()
	//截止此处，所有生成的消息已经发送了，下面是对数据库等内容进行的更新操作
	// 记录日志
	zlog.Info("success",
		zap.String("模型类型", model.Config.Model),
		zap.String("回复对象", sender.ID()),
		zap.String("记忆内容", memory.Memory),
		zap.String("消息内容", msg.Content),
		zap.String("回复内容", result),
	)

	// 更新 memory
	memory.MsgList = append(memory.MsgList, map[string]string{
		"role":    "user",
		"content": msg.Content,
	})
	memory.MsgList = append(memory.MsgList, map[string]string{
		"role":    "assistant",
		"content": result,
	})
	// 检查 MsgList 的长度和轮数
	const maxChars = 8888
	const maxRounds = 20
	totalChars := 0
	for _, entry := range memory.MsgList {
		totalChars += len(entry["content"])
	}

	// 如果超过限制，生成新的记忆，并删除旧的对话
	if totalChars > maxChars || len(memory.MsgList) > maxRounds*2 {
		// 使用大模型生成新的记忆
		newMemory := gen.GenMemory()
		// 更新 memory
		memory.Memory = newMemory
		if len(memory.MsgList) > maxRounds {
			memory.MsgList = memory.MsgList[len(memory.MsgList)/2:]
		}
	}
	update := bson.M{
		"$set": bson.M{
			"memory":  memory.Memory,
			"msglist": memory.MsgList,
		},
	}
	_, err = memoryctl.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		zlog.Error("更新 memory 和 msglist 失败", zap.Error(err))
	}
	zlog.Info("更新 memory 和 msglist 成功", zap.String("ID", sender.ID()), zap.String("更新后记忆:", memory.Memory))
}
