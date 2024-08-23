package auth

import (
	"1aides/internal/friends"
	"1aides/internal/groups"
	"1aides/pkg/log/zlog"

	"github.com/eatmoreapple/openwechat"
	"go.uber.org/zap"
)

func Auth(msg *openwechat.Message) bool {
	sender, err := msg.Sender()
	if err != nil {
		return false
	}

	if sender.IsSelf() {
		return false
	}

	var hasPermission bool
	if sender.IsGroup() {
		hasPermission, err = groups.CheckPermission(sender.ID())
	} else {
		hasPermission, err = friends.CheckPermission(sender.ID())
	}
	if err != nil {
		zlog.Debug("检查权限失败", zap.Error(err))
		return false
	}

	if !sender.IsGroup() && !hasPermission {
		msg.ReplyText("您没有权限使用此功能,请先申请激活码进行激活")
	}

	return hasPermission
}
