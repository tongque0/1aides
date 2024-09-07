package helper

import (
	"1aides/internal/friends"
	"1aides/internal/groups"
	"strings"

	"github.com/eatmoreapple/openwechat"
)

func setFriendMemory(command string, msg *openwechat.Message) {
	sender, err := msg.Sender()
	if err != nil {
		msg.ReplyText("获取发送者信息失败")
	}

	if sender.IsGroup() {
		msg.ReplyText("暂不支持群聊主动设置记忆，群聊记忆将自动更新")
		return
	}
	memoryContent := strings.TrimSpace(command[len("aides set memory"):])

	err = friends.SetFriendMemory(sender.ID(), memoryContent)
	if err != nil {
		msg.ReplyText("记忆内容更新失败")
	}

	msg.ReplyText("记忆内容已更新:" + memoryContent)
}

func getMemory(msg *openwechat.Message) {
	sender, err := msg.Sender()
	if err != nil {
		msg.ReplyText("获取发送者信息失败")
	}

	var memoryContent string

	if sender.IsGroup() {
		groupDetail, err := groups.GetGroupDetail(sender.ID())
		if err != nil {
			msg.ReplyText("获取群聊信息失败")
		}
		memoryContent = groupDetail.Memory
	} else {
		friendDetail, err := friends.GetFriendDetail(sender.ID())
		if err != nil {
			msg.ReplyText("获取好友信息失败")
		}
		memoryContent = friendDetail.Memory
	}

	msg.ReplyText("当前记忆为:\n" + memoryContent)
}
