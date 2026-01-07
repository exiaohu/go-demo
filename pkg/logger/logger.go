package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

// Initialize 初始化日志系统
func Initialize(debug bool) error {
	var config zap.Config
	if debug {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}

	// 配置日志输出格式
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	// 创建 logger
	var err error
	log, err = config.Build()
	if err != nil {
		return err
	}

	// 替换全局 logger
	zap.ReplaceGlobals(log)

	return nil
}

// Info 记录信息级别的日志
func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

// Warn 记录警告级别的日志
func Warn(msg string, fields ...zap.Field) {
	log.Warn(msg, fields...)
}

// Error 记录错误级别的日志
func Error(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
}

// Debug 记录调试级别的日志
func Debug(msg string, fields ...zap.Field) {
	log.Debug(msg, fields...)
}

// Fatal 记录致命级别的日志并退出程序
func Fatal(msg string, fields ...zap.Field) {
	log.Fatal(msg, fields...)
}

// Sync 刷新缓冲区
func Sync() error {
	if log != nil {
		return log.Sync()
	}
	return nil
}
