package logs

import (
	"testing"
	"time"
)

func TestFile(t *testing.T) {
	t.Log("测试fileLogWriter")
	config := `{
		"filename": "logs/test.log",
		"level" : 7
	}`
	fw := NewFileWriter()
	fw.Init(config)
	defer fw.Destroy()

	when := time.Now()
	err := fw.WriteMsg(when, "测试消息, 存入文件", 1)
	if err != nil {
		t.Errorf("写入错误, err: %v \n", err)
	}
}

func TestFileLog(t *testing.T) {
	t.Log("测试fileLog")

	config := `{
		"filename": "logs/test.log",
		"level" : 7
	}`

	logger := NewLogger()
	err := logger.SetLogger(AdapterFile, config)
	if err != nil {
		t.Errorf("写入错误, err: %v \n", err)
	}
	logger.Info("测试文件消息")
	logger.Close()
}
