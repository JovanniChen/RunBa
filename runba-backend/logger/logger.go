package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// LogLevel 日志级别
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

// String 返回日志级别的字符串表示
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// ParseLogLevel 将字符串转换为日志级别
func ParseLogLevel(level string) LogLevel {
	switch strings.ToLower(level) {
	case "debug":
		return DEBUG
	case "info":
		return INFO
	case "warn":
		return WARN
	case "error":
		return ERROR
	case "fatal":
		return FATAL
	default:
		return INFO // 默认为INFO级别
	}
}

// Color 返回日志级别对应的颜色代码
func (l LogLevel) Color() string {
	switch l {
	case DEBUG:
		return "\033[36m" // Cyan
	case INFO:
		return "\033[32m" // Green
	case WARN:
		return "\033[33m" // Yellow
	case ERROR:
		return "\033[31m" // Red
	case FATAL:
		return "\033[35m" // Magenta
	default:
		return "\033[0m" // Reset
	}
}

// Logger 改进的日志记录器
type Logger struct {
	level      LogLevel
	consoleLog *log.Logger
	fileLog    *log.Logger
	logFile    *os.File
	module     string
}

// Config 日志配置
type Config struct {
	Level        LogLevel // 日志级别
	LogDir       string   // 日志目录
	LogFile      string   // 日志文件名
	MaxFileSize  int64    // 最大文件大小（字节）
	EnableColor  bool     // 是否启用控制台颜色
	EnableCaller bool     // 是否显示调用位置
}

// DefaultConfig 默认配置
var DefaultConfig = Config{
	Level:        INFO,
	LogDir:       "logs",
	LogFile:      "app.log",
	MaxFileSize:  100 * 1024 * 1024, // 100MB
	EnableColor:  true,
	EnableCaller: true,
}

// globalLogger 全局日志实例
var globalLogger *Logger

// NewLogger 创建新的日志实例
func NewLogger(module string, config Config) *Logger {
	// 确保日志目录存在
	if err := os.MkdirAll(config.LogDir, 0755); err != nil {
		log.Printf("创建日志目录失败: %v", err)
		return nil
	}

	// 打开日志文件
	logFilePath := filepath.Join(config.LogDir, config.LogFile)
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("打开日志文件失败: %v", err)
		return nil
	}

	// 创建多重写入器（同时写入控制台和文件）
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	logger := &Logger{
		level:      config.Level,
		consoleLog: log.New(os.Stdout, "", 0),
		fileLog:    log.New(multiWriter, "", 0),
		logFile:    logFile,
		module:     module,
	}

	return logger
}

// InitGlobalLogger 初始化全局日志
func InitGlobalLogger(module string, config Config) {
	globalLogger = NewLogger(module, config)
	if globalLogger == nil {
		log.Fatal("初始化全局日志失败")
	}
}

// GetGlobalLogger 获取全局日志实例
func GetGlobalLogger() *Logger {
	if globalLogger == nil {
		// 使用默认配置初始化
		InitGlobalLogger("steam-backend", DefaultConfig)
	}
	return globalLogger
}

// Close 关闭日志文件
func (l *Logger) Close() error {
	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}

// formatMessage 格式化日志消息
func (l *Logger) formatMessage(level LogLevel, format string, args ...interface{}) (string, string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")

	// 获取调用者信息
	caller := ""
	if DefaultConfig.EnableCaller {
		_, file, line, ok := runtime.Caller(3) // 跳过3层调用栈
		if ok {
			filename := filepath.Base(file)
			caller = fmt.Sprintf(" [%s:%d]", filename, line)
		}
	}

	// 格式化消息内容
	message := fmt.Sprintf(format, args...)
	if message == "" {
		message = fmt.Sprint(args...)
	}

	// 控制台格式（带颜色）
	consoleMsg := ""
	if DefaultConfig.EnableColor {
		consoleMsg = fmt.Sprintf("%s[%s]%s [%s]%s%s - %s",
			level.Color(), level.String(), "\033[0m",
			timestamp, caller,
			func() string {
				if l.module != "" {
					return fmt.Sprintf(" <%s>", l.module)
				}
				return ""
			}(),
			message)
	} else {
		consoleMsg = fmt.Sprintf("[%s] [%s]%s%s - %s",
			level.String(), timestamp, caller,
			func() string {
				if l.module != "" {
					return fmt.Sprintf(" <%s>", l.module)
				}
				return ""
			}(),
			message)
	}

	// 文件格式（无颜色）
	fileMsg := fmt.Sprintf("[%s] [%s]%s%s - %s\n",
		level.String(), timestamp, caller,
		func() string {
			if l.module != "" {
				return fmt.Sprintf(" <%s>", l.module)
			}
			return ""
		}(),
		message)

	return consoleMsg, fileMsg
}

// log 内部日志方法
func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	if l == nil || level < l.level {
		return
	}

	consoleMsg, fileMsg := l.formatMessage(level, format, args...)

	// 输出到控制台
	fmt.Println(consoleMsg)

	// 输出到文件
	if l.logFile != nil {
		l.logFile.WriteString(fileMsg)
	}
}

// Debug 调试级别日志
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

// Info 信息级别日志
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

// Warn 警告级别日志
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WARN, format, args...)
}

// Error 错误级别日志
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

// Fatal 致命错误级别日志
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(FATAL, format, args...)
	os.Exit(1)
}

// 全局便捷方法
func Debug(format string, args ...interface{}) {
	GetGlobalLogger().Debug(format, args...)
}

func Info(format string, args ...interface{}) {
	GetGlobalLogger().Info(format, args...)
}

func Warn(format string, args ...interface{}) {
	GetGlobalLogger().Warn(format, args...)
}

func Error(format string, args ...interface{}) {
	GetGlobalLogger().Error(format, args...)
}

func Fatal(format string, args ...interface{}) {
	GetGlobalLogger().Fatal(format, args...)
}

// WithFields 带字段的日志（用于结构化日志）
func (l *Logger) WithFields(fields map[string]interface{}) *FieldLogger {
	return &FieldLogger{
		logger: l,
		fields: fields,
	}
}

// FieldLogger 带字段的日志记录器
type FieldLogger struct {
	logger *Logger
	fields map[string]interface{}
}

// formatFields 格式化字段
func (fl *FieldLogger) formatFields() string {
	if len(fl.fields) == 0 {
		return ""
	}

	var parts []string
	for key, value := range fl.fields {
		parts = append(parts, fmt.Sprintf("%s=%v", key, value))
	}
	return fmt.Sprintf(" {%s}", strings.Join(parts, ", "))
}

// Debug 带字段的调试日志
func (fl *FieldLogger) Debug(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...) + fl.formatFields()
	fl.logger.Debug("%s", message)
}

// Info 带字段的信息日志
func (fl *FieldLogger) Info(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...) + fl.formatFields()
	fl.logger.Info("%s", message)
}

// Warn 带字段的警告日志
func (fl *FieldLogger) Warn(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...) + fl.formatFields()
	fl.logger.Warn("%s", message)
}

// Error 带字段的错误日志
func (fl *FieldLogger) Error(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...) + fl.formatFields()
	fl.logger.Error("%s", message)
}

// Fatal 带字段的致命错误日志
func (fl *FieldLogger) Fatal(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...) + fl.formatFields()
	fl.logger.Fatal("%s", message)
}
