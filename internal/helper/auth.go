package helper

import (
	"1aides/internal/friends"
	"1aides/internal/groups"
	"1aides/pkg/log/zlog"
	"strings"

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
	// 如果是群聊，检查消息是否包含特定前缀或关键字
	if sender.IsGroup() && hasPermission {
		content := msg.Content
		if !strings.HasPrefix(content, "@小喜") && !strings.Contains(content, "小喜") {
			return false
		}
	}

	if !sender.IsGroup() && !hasPermission {
		msg.ReplyText("您没有权限使用此功能,请先申请激活码进行激活")
	}

	return hasPermission
}
