package zlog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

// init 在包初始化时配置和构建全局的日志器。
func init() {
	// 配置日志格式和级别
	config := zap.NewProductionConfig()
	// 设置时间戳字段名称为 "timestamp"
	config.EncoderConfig.TimeKey = "timestamp"
	// 使用ISO8601标准格式化时间戳
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// 构建配置的日志器
	var err error
	logger, err = config.Build(zap.AddCallerSkip(1))
	if err != nil {
		// 初始化日志器失败时，抛出panic异常
		panic(err)
	}
}

// Info 记录一条信息级别的日志。
func Info(message string, fields ...zap.Field) {
	logger.Info(message, fields...)
}

// Error 记录一条错误级别的日志。
func Error(message string, fields ...zap.Field) {
	logger.Error(message, fields...)
}

// Debug 记录一条调试级别的日志。
func Debug(message string, fields ...zap.Field) {
	logger.Debug(message, fields...)
}

// Warn 记录一条警告级别的日志。
func Warn(message string, fields ...zap.Field) {
	logger.Warn(message, fields...)
}

// Fatal 记录一条致命错误级别的日志，并终止程序。
func Fatal(message string, fields ...zap.Field) {
	logger.Fatal(message, fields...)
}

// Sync 刷新所有缓冲的日志条目。
// 应在程序关闭前调用，以确保所有日志都被写入目标位置。
func Sync() error {
	return logger.Sync()
}
