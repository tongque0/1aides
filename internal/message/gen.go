package message

import (
	"1aides/pkg/common/config"
	"1aides/pkg/generator"
	"1aides/pkg/generator/modhub"
	"1aides/pkg/log/zlog"
	"os"

	"github.com/eatmoreapple/openwechat"
	"go.uber.org/zap"
)

func gen(msg *openwechat.Message) {

	consulAddress := os.Getenv("CONSUL_ADDRESS")
	if consulAddress == "" {
		consulAddress = "127.0.0.1:8500"
	}

	cfg, err := config.NewConsulConfig(consulAddress)
	if err != nil {
		zlog.Error("创建配置管理器失败", zap.Any("error", err))
		return
	}

	// 从Consul加载配置
	var model modhub.Model
	err = cfg.LoadConfig("1aides/model", &model)
	if err != nil {
		zlog.Error("加载配置失败", zap.Any("error", err))
		return
	}
	zap.S().Infof("Loaded model config: %+v", model)

	// 使用加载的模型配置创建生成器实例
	gen := generator.NewGenerator(msg, generator.WithModel(model))

	// 生成响应
	gen.Generate()
}
