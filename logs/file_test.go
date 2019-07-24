package logs

import (
	"testing"
	"time"
)

func TestFileMode(t *testing.T) {
	t.Log("测试fileLogWriter")
	config := `{
		"filename": "logs/test.log",
		"level" : 7,
		"maxlines": 100
	}`
	fw := NewFileWriter()
	fw.Init(config)
	defer fw.Destroy()

	when := time.Now()
	var err error
	for i := 0; i < 101; i++ {
		err = fw.WriteMsg(when, "测试消息, 存入文件", 1)
		if err != nil {
			t.Errorf("写入错误, err: %v \n", err)
		}
	}

}

func TestFileMuti(t *testing.T) {
	t.Log("多 logger 测试")
	config := `{
		"filename": "logs/test.log"
	}`

	log1 := NewLogger()
	log1.SetLogger(AdapterFile, config)
	log1.Info("test1")

	log2 := NewLogger()
	log2.SetLogger(AdapterFile, config)
	log2.Info("test2")
}

func TestFileLog(t *testing.T) {
	t.Log("测试 fileLog")

	config := `{
		"filename": "logs/test.log",
		"level" : 7
	}`

	logger := NewLogger()
	err := logger.SetLogger(AdapterFile, config)
	if err != nil {
		t.Errorf("写入错误, err: %v \n", err)
	}

	msg := "测试文件消息"
	logger.Emergency(msg)
	logger.Alert(msg)
	logger.Critical(msg)
	logger.Error(msg)
	logger.Warning(msg)
	logger.Warn(msg)
	logger.Notice(msg)
	logger.Info(msg)
	logger.Debug(msg)
	logger.Close()
}

// 同时命令行和文件输出
func TestMultiMode(t *testing.T) {
	t.Log("测试 fileLog")

	config := `{
		"filename": "logs/test.log",
		"level" : 7
	}`

	// 输出文件模式
	err := SetLogger(AdapterFile, config)
	if err != nil {
		t.Errorf("写入错误, err: %v \n", err)
	}

	// 开启打印命令行模式
	err = SetLogger(AdapterConsole)
	if err != nil {
		t.Errorf("写入错误, err: %v \n", err)
	}

	msg := "测试文件消息"
	Emergency(msg)
	Alert(msg)
	Critical(msg)
	Error(msg)
	Warning(msg)
	Warn(msg)
	Notice(msg)
	Info(msg)
	Debug(msg)
}
