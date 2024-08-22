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
	group, err := self.Groups()
	zlog.Info("获取所有的群组", zap.Any("群组", group), zap.Error(err))
}
