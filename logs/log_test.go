package logs

import (
	"testing"
)

func TestLog(t *testing.T) {
	callAll(LevelDebug, "测试消息")
}

func callAll(level int, msg string) {
	logger := NewLogger()
	logger.SetLevel(level)

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
