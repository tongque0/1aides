package friends

import (
	"1aides/pkg/components/bot"
	"1aides/pkg/log/zlog"

	"go.uber.org/zap"
)

func GetFriends() {
	// 获取所有的好友
	self, err := bot.WxBot.GetCurrentUser()
	if err != nil {
		zlog.Error("获取当前用户失败", zap.Error(err))
		return
	}
	friends, err := self.Friends()
	zlog.Info("获取所有的好友", zap.Any("friends", friends), zap.Error(err))
}
