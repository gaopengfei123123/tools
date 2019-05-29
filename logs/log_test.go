package logs

import (
	"fmt"
	"testing"
)

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
	for _, level := range levelPool {
		callAll(level, "测试消息")
	}

}

func callAll(level int, msg string) {
	logger := NewLogger()
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
