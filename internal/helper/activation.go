package helper

import (
	"1aides/pkg/components/db"
	"1aides/pkg/log/zlog"
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/eatmoreapple/openwechat"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// Activation 结构体表示验证码及其验证状态
type Activation struct {
	Code       string    `bson:"code"`
	IsVerified bool      `bson:"is_verified"`
	CreatedAt  time.Time `bson:"created_at"`
}

// GenerateActivationCode 生成随机验证码并保存到数据库
func GenerateActivationCode() string {
	collection := db.GetMongoDB().Collection("activation")

	// 生成随机验证码
	code := generateRandomCode()

	activation := Activation{
		Code:       code,
		IsVerified: false,
		CreatedAt:  time.Now(),
	}

	_, err := collection.InsertOne(context.TODO(), activation)
	if err != nil {
		zlog.Error("插入验证码失败", zap.Error(err))
		return ""
	}

	return code
}

// VerifyActivationCode 验证验证码
func VerifyActivationCode(command string, msg *openwechat.Message) bool {
	sender, err := msg.Sender()
	if err != nil {
		zlog.Error("Failed to get sender:", zap.Error(err))
		return false
	}
	// 从命令字符串中提取验证码
	parts := strings.Fields(command)
	if len(parts) < 3 {
		zlog.Error("Invalid command format")
		return false
	}
	code := parts[2]

	collection := db.GetMongoDB().Collection("activation")

	filter := bson.M{"code": code, "is_verified": false}

	var activation Activation
	err = collection.FindOne(context.TODO(), filter).Decode(&activation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
		zlog.Error("Failed to find activation code:", zap.Error(err))
	}

	// 更新验证码状态为已验证
	_, err = collection.UpdateOne(
		context.TODO(),
		filter,
		bson.M{"$set": bson.M{"is_verified": true}},
	)
	if err != nil {
		zlog.Error("Failed to update activation code:", zap.Error(err))
	}

	// 验证码验证通过后，更新用户权限
	err = updateUserPermission(sender.ID())
	if err != nil {
		zlog.Error("Failed to update user permission:", zap.Error(err))
	}

	return true
}

// updateUserPermission 更新用户权限
func updateUserPermission(senderID string) error {
	// 定义MongoDB集合
	friendsCollection := db.GetMongoDB().Collection("friends")
	groupsCollection := db.GetMongoDB().Collection("groups")

	// 在 friends 表中查找用户并更新权限
	filter := bson.M{"id": senderID}
	update := bson.M{"$set": bson.M{"has_permission": true}}

	// 首先尝试在 friends 集合中更新
	result, err := friendsCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	// 如果 matchedCount 大于 0，表示在 friends 集合中找到了并更新了用户
	if result.MatchedCount > 0 {
		return nil
	}

	// 如果在 friends 集合中未找到用户，则尝试在 groups 集合中更新
	result, err = groupsCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	// 如果在 groups 集合中找到了并更新了用户，返回 nil
	if result.MatchedCount > 0 {
		return nil
	}

	// 如果两个集合中都没有找到用户，返回错误
	return fmt.Errorf("user with ID %s not found in friends or groups", senderID)
}

func generateRandomCode() string {
	// 获取当前时间的时间戳（精确到毫秒）
	timestamp := time.Now().UnixNano() / 1e6

	// 生成一个随机数，确保ID的唯一性
	randomPart := rand.Int63() & 0xFFFFFF // 生成24位的随机数部分

	// 将时间戳和随机部分拼接成一个16进制字符串
	id := fmt.Sprintf("%x%x", timestamp, randomPart)

	return id
}
