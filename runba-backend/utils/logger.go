package utils

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

// InitLogger 初始化日志系统
func InitLogger() {
	// 配置日志编码器
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// 根据环境设置日志级别
	var level zapcore.Level
	if os.Getenv("GIN_MODE") == "release" {
		level = zapcore.InfoLevel
	} else {
		level = zapcore.DebugLevel
	}

	// 创建核心配置
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		level,
	)

	// 创建logger
	logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}

// GetLogger 获取logger实例
func GetLogger() *zap.Logger {
	if logger == nil {
		InitLogger()
	}
	return logger
}

// LogRequest 记录HTTP请求日志
func LogRequest(method, path string, statusCode int, latency time.Duration, clientIP string) {
	GetLogger().Info("HTTP请求",
		zap.String("method", method),
		zap.String("path", path),
		zap.Int("status_code", statusCode),
		zap.Duration("latency", latency),
		zap.String("client_ip", clientIP),
	)
}

// LogError 记录错误日志
func LogError(message string, err error, fields ...zap.Field) {
	GetLogger().Error(message,
		append(fields, zap.Error(err))...,
	)
}

// LogInfo 记录信息日志
func LogInfo(message string, fields ...zap.Field) {
	GetLogger().Info(message, fields...)
}

// LogDebug 记录调试日志
func LogDebug(message string, fields ...zap.Field) {
	GetLogger().Debug(message, fields...)
}

// LogWarn 记录警告日志
func LogWarn(message string, fields ...zap.Field) {
	GetLogger().Warn(message, fields...)
}
