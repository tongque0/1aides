package helper

import (
	"1aides/internal/friends"
	"regexp"
	"strings"

	"github.com/eatmoreapple/openwechat"
)

// 定义命令行帮助信息为常量字符串
const CommandList = `NAME:
   1aides - 可靠敏捷的微信机器人服务

USAGE:
   aides [command]

COMMANDS:
   help        帮助手册

   active 	   激活服务
		Example: aides active [激活码]

   close       暂停服务

   set		   设置服务

		Example: aides set memory [记忆内容]

   get		   获取服务信息

		Example: aides get memory
				 aides get memory

GLOBAL OPTIONS:
   --help, -h  Show help
`

const adminlist = `
COMMANDS:
    help        帮助手册

	gencode    生成激活码

`

// IsHelp 检查消息是否是帮助命令
func IsHelp(msg *openwechat.Message) bool {
	match, _ := regexp.MatchString(`^aides`, strings.ToLower(msg.Content))
	return match
}

func GroupRule(msg *openwechat.Message) bool {

	return false
}

// Helper 处理匹配到的命令
func Helper(msg *openwechat.Message) {
	// 获取命令内容
	command := strings.TrimSpace(msg.Content)
	switch {
	case strings.HasPrefix(command, "aides help"):
		msg.ReplyText(CommandList)
	case strings.HasPrefix(command, "aides active"):
		VerifyActivationCode(command, msg)
	case strings.HasPrefix(command, "aides close"):
		msg.ReplyText("成功关闭服务")
	case strings.HasPrefix(command, "aides set"):
		setCmd(command, msg)
	case strings.HasPrefix(command, "aides get"):
		getCmd(command, msg)
	case strings.HasPrefix(command, "aides admin"):
		adminCmd(command, msg)
	default:
		msg.ReplyText(CommandList)
	}
}

func adminCmd(command string, msg *openwechat.Message) {
	sender, err := msg.Sender()
	if err != nil {
		return
	}
	isAdmin, err := friends.CheckAdmin(sender.ID())
	if err != nil {
		msg.ReplyText("您没有权限使用此功能")
		return
	}
	if !isAdmin {
		msg.ReplyText("您没有权限使用此功能")
		return
	}
	switch {
	case strings.HasPrefix(command, "aides admin help"):
		msg.ReplyText(adminlist)
	case strings.HasPrefix(command, "aides admin gencode"):
		GenerateActivationCode(command, msg)
	default:
		msg.ReplyText(adminlist)
	}
}

func setCmd(command string, msg *openwechat.Message) {
	switch {
	case strings.HasPrefix(command, "aides set memory"):
		setFriendMemory(command, msg)
	default:
		msg.ReplyText("服务当前状态：运行中")
	}
}

func getCmd(command string, msg *openwechat.Message) {
	switch {
	case strings.HasPrefix(command, "aides get memory"):
		getMemory(msg)
	case strings.HasPrefix(command, "aides get history"):
		getHistory(msg)
	default:
		msg.ReplyText("服务当前状态：运行中")
	}
}
