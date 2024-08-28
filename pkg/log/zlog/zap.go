package zlog

import (
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

// init 在包初始化时配置和构建全局的日志器。
func init() {
	// 配置日志格式和级别
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// 文件日志配置
	logFile := &lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    100,
		MaxBackups: 7,
		MaxAge:     30,
		Compress:   true,
	}

	// 文件日志的Core
	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(config.EncoderConfig), // JSON编码器
		zapcore.AddSync(logFile),                     // 日志切割
		zapcore.InfoLevel,                            // 文件日志级别
	)

	// 控制台输出配置
	consoleEncoderConfig := zap.NewDevelopmentEncoderConfig()
	consoleEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // 控制台彩色输出

	// 控制台日志的Core
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(consoleEncoderConfig), // 控制台编码器
		zapcore.AddSync(os.Stdout),                      // 标准输出
		zapcore.DebugLevel,                              // 控制台日志级别
	)

	// 使用MultiCore组合多个日志Core
	core := zapcore.NewTee(
		fileCore,
		consoleCore,
	)

	// 创建日志器并应用核心
	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
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
