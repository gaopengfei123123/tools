package logs

import (
	"fmt"
	"testing"
)

func TestFormater(t *testing.T) {
	var formatter FormaterFunc
	formatter = func(logLevel int, msg string, v ...interface{}) string {
		return fmt.Sprintf("LogFormatter -> logLevel: %d, msg: %s, other: %v \n", logLevel, msg, v)
	}

	logger := NewLogger()
	logger.SetFormatter(formatter)

	logger.Info("这是一条测试消息", 233, 244, 88)
}

func TestLog(t *testing.T) {
	levelPool := []int{
		LevelEmergency,
		LevelAlert,
		LevelCritical,
		LevelError,
		LevelWarning,
		LevelNotice,
		LevelInformational,
		LevelDebug,
	}
	logger := NewLogger()
	for _, level := range levelPool {
		callAll(level, "测试消息", logger)
	}

}

func callAll(level int, msg string, logger *Logger) {
	logger.SetLevel(level)
	fmt.Printf(">>>>>>>>>>>>>>>>>>>>>>>>level: %d \n", level)

	logger.Emergency(msg)
	logger.Alert(msg)
	logger.Critical(msg)
	logger.Error(msg)
	logger.Warning(msg)
	logger.Warn(msg)
	logger.Notice(msg)
	logger.Info(msg)
	logger.Debug(msg)
}
