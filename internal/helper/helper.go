package helper

import (
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
   help        Show this help message

   open        Activate the service

   close       Deactivate the service

   status      Show the status of the service

   restart     Restart the service


GLOBAL OPTIONS:
   --help, -h  Show help
`

const adminlist = `
COMMANDS:
	gencode    生成激活码

`

// IsHelp 检查消息是否是帮助命令
func IsHelp(msg *openwechat.Message) bool {
	match, _ := regexp.MatchString(`^aides`, strings.ToLower(msg.Content))
	return match
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
	case strings.HasPrefix(command, "aides status"):
		msg.ReplyText("服务当前状态：运行中")
	case strings.HasPrefix(command, "aides restart"):
		adminCmd(command, msg)
	case strings.HasPrefix(command, "aides admin"):
		adminCmd(command, msg)
	default:
		msg.ReplyText(CommandList)
	}
}

func adminCmd(command string, msg *openwechat.Message) {
	switch {
	case strings.HasPrefix(command, "aides admin help"):
		msg.ReplyText(adminlist)
	case strings.HasPrefix(command, "aides admin gencode"):
		msg.ReplyText(GenerateActivationCode())
	default:
		msg.ReplyText(adminlist)
	}
}
