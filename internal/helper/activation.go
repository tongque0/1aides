package helper

import (
	"1aides/pkg/components/db"
	"1aides/pkg/log/zlog"
	"context"
	"fmt"
	"math/rand"
	"strconv"
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
func GenerateActivationCode(command string, msg *openwechat.Message) {
	// 从命令字符串中提取参数
	parts := strings.Fields(command)
	var nums int

	// 检查是否有指定数量的参数
	if len(parts) >= 4 {
		// 将字符串转换为整数
		num, err := strconv.Atoi(parts[3])
		if err != nil {
			zlog.Error("无效的数量参数", zap.Error(err))
		}
		nums = num
	} else {
		nums = 1
	}

	collection := db.GetMongoDB().Collection("activation")
	var codes []string

	for i := 0; i < nums; i++ {
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
		}

		codes = append(codes, code)
	}
	msg.ReplyText("成功生成激活码:\n" + strings.Join(codes, "\n"))
}

// VerifyActivationCode 验证验证码
func VerifyActivationCode(command string, msg *openwechat.Message) {
	sender, err := msg.Sender()
	if err != nil {
		zlog.Error("Failed to get sender:", zap.Error(err))
	}
	// 从命令字符串中提取验证码
	parts := strings.Fields(command)
	if len(parts) < 3 {
		zlog.Error("Invalid command format")
	}
	code := parts[2]

	collection := db.GetMongoDB().Collection("activation")

	filter := bson.M{"code": code, "is_verified": false}

	var activation Activation
	err = collection.FindOne(context.TODO(), filter).Decode(&activation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
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

	msg.ReplyText("激活成功！")
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

	// 将时间戳转换为16进制字符串，并截取后4位
	timestampHex := fmt.Sprintf("%x", timestamp)
	truncatedTimestamp := timestampHex[len(timestampHex)-4:]

	// 生成一个随机数，确保ID的唯一性
	randomPart := rand.Int63() & 0xFFFFFF // 生成24位的随机数部分

	// 将固定前缀、截取后的时间戳和随机部分拼接成一个16进制字符串
	id := fmt.Sprintf("1aides%s%x", truncatedTimestamp, randomPart)

	return id
}
