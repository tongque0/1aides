package zlog

import (
	"github.com/natefinch/lumberjack"
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

	// 配置日志切割
	logFile := &lumberjack.Logger{
		Filename:   "logs/app.log", // 日志文件路径
		MaxSize:    100,            // 以MB为单位，日志文件达到该大小后切割
		MaxBackups: 7,              // 保留的旧日志文件最大数量
		MaxAge:     30,             // 日志文件最大保存天数
		Compress:   true,           // 是否压缩旧日志文件
	}

	// 创建日志核心（Core），指定日志级别、编码器和输出
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config.EncoderConfig), // 使用JSON格式记录日志
		zapcore.AddSync(logFile),                     // 日志写入到文件并支持切割
		zapcore.InfoLevel,                            // 最低日志记录级别
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
