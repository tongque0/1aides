package helper

import (
	"1aides/internal/friends"
	"1aides/internal/groups"

	"github.com/eatmoreapple/openwechat"
)

func getHistory(msg *openwechat.Message) {
	sender, err := msg.Sender()
	if err != nil {
		msg.ReplyText("获取发送者信息失败")
		return // Ensure to return after sending error message
	}

	var memoryContent string

	if sender.IsGroup() {
		groupDetail, err := groups.GetGroupDetail(sender.ID())
		if err != nil {
			msg.ReplyText("获取群聊信息失败")
			return
		}
		for _, msg := range groupDetail.MsgList {
			for key, value := range msg {
				memoryContent += key + ": " + value + "\n"
			}
		}
	} else {
		friendDetail, err := friends.GetFriendDetail(sender.ID())
		if err != nil {
			msg.ReplyText("获取好友信息失败")
			return
		}
		for _, msg := range friendDetail.MsgList {
			for key, value := range msg {
				memoryContent += key + ": " + value + "\n"
			}
		}
	}

	if memoryContent == "" {
		memoryContent = "无可用记忆"
	}

	//微信支持长度有限，这里可能存在bug
	msg.ReplyText("以下为聊天记录:\n" + memoryContent)
}
