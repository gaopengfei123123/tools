package logs

import (
	"fmt"
	"strings"
)

// 直接暴露在外的公用方法

// 默认全局的logger
var logger = NewLogger()

// GetLogger 给出去
func GetLogger() *Logger {
	return logger
}

// Reset 重置
func Reset() {
	logger.Reset()
}

// Async 设置异步长度
func Async(msgLen ...int64) *Logger {
	return logger.Async(msgLen...)
}

// SetLevel 设置日志等级
func SetLevel(l int) {
	logger.SetLevel(l)
}

// EnableFuncCallDepth 是否显示调用文件位置
func EnableFuncCallDepth(b bool) {
	logger.enableFuncCallDepth = b
}

// SetLogFuncCallDepth 设置文件调用层级
func SetLogFuncCallDepth(d int) {
	logger.loggerFuncCallDepth = d
}

// SetLogger 配置
func SetLogger(adapter string, config ...string) error {
	return logger.SetLogger(adapter, config...)
}

// Emergency logs a message at emergency level.
func Emergency(f interface{}, v ...interface{}) {
	logger.Emergency(formatLog(f, v...))
}

// Alert logs a message at alert level.
func Alert(f interface{}, v ...interface{}) {
	logger.Alert(formatLog(f, v...))
}

// Critical logs a message at critical level.
func Critical(f interface{}, v ...interface{}) {
	logger.Critical(formatLog(f, v...))
}

// Error logs a message at error level.
func Error(f interface{}, v ...interface{}) {
	logger.Error(formatLog(f, v...))
}

// Warning logs a message at warning level.
func Warning(f interface{}, v ...interface{}) {
	logger.Warn(formatLog(f, v...))
}

// Warn compatibility alias for Warning()
func Warn(f interface{}, v ...interface{}) {
	logger.Warn(formatLog(f, v...))
}

// Notice logs a message at notice level.
func Notice(f interface{}, v ...interface{}) {
	logger.Notice(formatLog(f, v...))
}

// Informational logs a message at info level.
func Informational(f interface{}, v ...interface{}) {
	logger.Info(formatLog(f, v...))
}

// Info compatibility alias for Warning()
func Info(f interface{}, v ...interface{}) {
	logger.Info(formatLog(f, v...))
}

// Debug logs a message at debug level.
func Debug(f interface{}, v ...interface{}) {
	logger.Debug(formatLog(f, v...))
}

// Trace logs a message at trace level.
// compatibility alias for Warning()
func Trace(f interface{}, v ...interface{}) {
	logger.Debug(formatLog(f, v...))
}

func formatLog(f interface{}, v ...interface{}) string {
	var msg string
	switch f.(type) {
	case string:
		msg = f.(string)
		if len(v) == 0 {
			return msg
		}
		if strings.Contains(msg, "%") && !strings.Contains(msg, "%%") {
			//format string
		} else {
			//do not contain format char
			msg += strings.Repeat(" %v", len(v))
		}
	default:
		msg = fmt.Sprint(f)
		if len(v) == 0 {
			return msg
		}
		msg += strings.Repeat(" %v", len(v))
	}
	return fmt.Sprintf(msg, v...)
}
