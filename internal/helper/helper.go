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
   help        帮助手册

   active 	   激活服务
		Example: aides active [激活码]

   close       暂停服务

   set		   设置服务


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
		GenerateActivationCode(command, msg)
	default:
		msg.ReplyText(adminlist)
	}
}

func setCmd(command string, msg *openwechat.Message) {
	switch {
	case strings.HasPrefix(command, "aides set prompt"):
		msg.ReplyText("服务当前状态：运行中")
	default:
		msg.ReplyText("服务当前状态：运行中")
	}
}
