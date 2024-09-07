package groups

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

type Group struct {
	ID            string              `bson:"id"`
	NickName      string              `bson:"nick_name"`
	RemarkName    string              `bson:"remark_name"`
	HasPermission bool                `bson:"has_permission"` // 权限字段
	Memory        string              `bson:"memory"`
	MsgList       []map[string]string `bson:"msglist"`
}

func InitGroupsDB() {
	// 获取当前用户
	self, err := bot.WxBot.GetCurrentUser()
	if err != nil {
		zlog.Error("获取当前用户失败", zap.Error(err))
		return
	}

	// 获取用户的所有群组
	groups, err := self.Groups()
	if err != nil {
		zlog.Error("获取群组列表失败", zap.Error(err))
		return
	}

	// 获取MongoDB集合
	collection := db.GetMongoDB().Collection("groups")

	// 用于存储当前所有群组的ID
	currentGroupIDs := make([]interface{}, 0, len(groups))

	// 插入或更新群组信息到数据库
	for _, group := range groups {
		currentGroupIDs = append(currentGroupIDs, group.ID()) // 添加当前群组ID到列表

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// MongoDB中使用filter来查找群组ID是否存在
		filter := bson.M{"id": group.ID()}

		// 使用$set和$setOnInsert
		update := bson.M{
			"$set": bson.M{
				"nick_name":   group.NickName,
				"remark_name": group.RemarkName,
			},
			"$setOnInsert": bson.M{
				"has_permission": false, // 新群组时初始化权限为false
			},
		}

		// 更新或者插入群组，如果群组存在则更新，不存在则插入
		opts := options.Update().SetUpsert(true)
		_, err := collection.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			zlog.Error("插入或更新群组信息失败", zap.Error(err))
		} else {
			zlog.Info("成功插入或更新群组信息", zap.String("ID", group.ID()))
		}
	}

	// 删除数据库中不在当前群组列表中的记录
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	deletionFilter := bson.M{"id": bson.M{"$nin": currentGroupIDs}}
	deleteResult, err := collection.DeleteMany(ctx, deletionFilter)
	if err != nil {
		zlog.Error("删除非群组记录失败", zap.Error(err))
	} else {
		zlog.Info("已删除非群组记录", zap.Int64("count", deleteResult.DeletedCount))
	}
}

// CheckPermission 检查指定好友是否具有权限
func CheckPermission(groupsID string) (bool, error) {
	collection := db.GetMongoDB().Collection("groups")
	filter := bson.M{"id": groupsID}

	var group Group
	err := collection.FindOne(context.Background(), filter).Decode(&group)
	if err != nil {
		zlog.Error("查询群组权限失败", zap.String("群组id", groupsID), zap.Error(err))
		return false, err
	}

	return group.HasPermission, nil
}

// GetGroups 获取所有群组
func GetGroups() []Group {
	collection := db.GetMongoDB().Collection("groups")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		zlog.Error("查询群组失败", zap.Error(err))
		return nil
	}
	defer cursor.Close(context.Background())

	var groups []Group
	for cursor.Next(context.Background()) {
		var group Group
		err := cursor.Decode(&group)
		if err != nil {
			zlog.Error("解码群组失败", zap.Error(err))
			continue
		}
		groups = append(groups, group)
	}
	return groups
}

// SetPermission 设置好友权限
func SetPermission(friendID string, permission bool) error {
	collection := db.GetMongoDB().Collection("groups")
	filter := bson.M{"id": friendID}
	update := bson.M{"$set": bson.M{"has_permission": permission}}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		zlog.Error("设置群组权限失败", zap.Error(err))
		return err
	}
	return nil
}

func GetGroupDetail(groupID string) (Group, error) {
	collection := db.GetMongoDB().Collection("groups")
	filter := bson.M{"id": groupID}

	var group Group
	err := collection.FindOne(context.Background(), filter).Decode(&group)
	if err != nil {
		zlog.Error("查询群组失败", zap.Error(err))
		return Group{}, err
	}
	return group, nil
}
