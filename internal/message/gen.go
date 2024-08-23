package message

import (
	"1aides/internal/friends"
	"1aides/pkg/common/config"
	"1aides/pkg/components/db"
	"1aides/pkg/components/generator"
	"1aides/pkg/components/generator/memory"
	"1aides/pkg/components/generator/modhub"
	"1aides/pkg/log/zlog"
	"context"
	"os"

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
	consulAddress := os.Getenv("CONSUL_ADDRESS")
	if consulAddress == "" {
		consulAddress = "127.0.0.1:8500"
	}

	cfg, err := config.NewConsulConfig(consulAddress)
	if err != nil {
		zlog.Error("创建配置管理器失败", zap.Error(err))
		return
	}

	// 从配置中心加载模型配置
	var model modhub.Model
	err = cfg.LoadConfig("1aides/model", &model)
	if err != nil {
		zlog.Error("加载配置失败", zap.Error(err))
		return
	}
	// 从mongoDB中加载记忆内容
	var memory memory.Memory
	var collection *mongo.Collection
	if sender.IsGroup() {
		collection = db.GetMongoDB().Collection("groups")
	} else {
		collection = db.GetMongoDB().Collection("friends")
	}
	filter := bson.M{"id": sender.ID()}
	var friend friends.Friend
	err = collection.FindOne(context.TODO(), filter).Decode(&friend)
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
	const maxChars = 4000
	const maxRounds = 10
	totalChars := 0
	for _, entry := range memory.MsgList {
		totalChars += len(entry["content"])
	}

	// 如果超过限制，生成新的记忆，并删除旧的对话
	if totalChars > maxChars || len(memory.MsgList) > maxRounds*2 {
		// 使用大模型生成新的记忆
		// newMemory := generateMemoryFromMsgList(memory.MsgList)

		// 更新 memory
		memory.Memory = "新的记忆"

		// 删除旧的 5 轮对话（10 条消息）
		if len(memory.MsgList) > 20 {
			memory.MsgList = memory.MsgList[10:]
		}
	}
	update := bson.M{
		"$set": bson.M{
			"memory":  memory.Memory,
			"msglist": memory.MsgList,
		},
	}
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		zlog.Error("更新 memory 和 msglist 失败", zap.Error(err))
	}
	zlog.Info("更新 memory 和 msglist 成功", zap.String("ID", sender.ID()))
}
