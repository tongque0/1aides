package friends

import (
	"1aides/pkg/components/bot"
	"1aides/pkg/components/db"
	"1aides/pkg/log/zlog"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// Friend 定义好友信息的结构体
type Friend struct {
	ID            string `bson:"id"`
	NickName      string `bson:"nick_name"`
	RemarkName    string `bson:"remark_name"`
	HasPermission bool   `bson:"has_permission"` // 权限字段
}

// InitFriendDB 初始化好友数据库
func InitFriendDB() {
	// 获取所有的好友
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

	// 获取MongoDB集合
	collection := db.GetMongoDB().Collection("friends")

	// 插入好友信息到数据库
	for _, friend := range friends {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// MongoDB中使用filter来查找用户ID是否存在
		filter := bson.M{"id": friend.ID()}

		// 使用$set和$setOnInsert
		update := bson.M{
			"$set": bson.M{
				"nick_name":   friend.NickName,
				"remark_name": friend.RemarkName,
			},
			"$setOnInsert": bson.M{
				"has_permission": false, // 新用户时初始化权限为false
			},
		}

		// 更新或者插入用户，如果用户存在则更新，不存在则插入
		opts := options.Update().SetUpsert(true)
		_, err := collection.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			zlog.Error("插入或更新好友信息失败", zap.Error(err))
		} else {
			zlog.Info("成功插入或更新好友信息", zap.String("ID", friend.ID()))
		}
	}
}

// CheckPermission 检查指定好友是否具有权限
func CheckPermission(friendID string) (bool, error) {
	collection := db.GetMongoDB().Collection("friends")
	filter := bson.M{"id": friendID}

	var friend Friend
	err := collection.FindOne(context.Background(), filter).Decode(&friend)
	if err != nil {
		zlog.Error("查询好友权限失败", zap.Error(err))
		return false, err
	}

	return friend.HasPermission, nil
}
