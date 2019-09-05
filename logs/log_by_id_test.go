package logs

import (
	// "fmt"
	"sync"
	"testing"
	// "time"
)

func TestLogDebug(t *testing.T) {
	config := `{
		"filename": "logs/test.log",
		"level" : 7,
		"maxlines": 1000000
	}`

	logger := NewLoggerByID()
	// defer logger.Close()

	logger.SetLogger(AdapterFile, config)
	uuid, _ := GenerageUniqueID("test")
	logger.ID = uuid
	logger.IP = "127.0.0.1"
	logger.Info("测试消息", "ERROR_DEFAULT")
}

func TestLogByID(t *testing.T) {
	t.Log("测试LogByID")
	config := `{
		"filename": "logs/test.log",
		"level" : 7,
		"maxlines": 1000000
	}`

	logger := NewLoggerByID()
	// defer logger.Close()

	logger.SetLogger(AdapterFile, config)

	// when := time.Now()
	// var err error
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(index int) {
			t.Logf("start log %d \n", index)
			logger := NewLoggerByID()
			// defer logger.Close()

			logger.SetLogger(AdapterFile, config)
			logger.SetLogger(AdapterConsole)

			uuid, _ := GenerageUniqueID("test")
			logger.ID = uuid
			logger.IP = "127.0.0.1"
			logger.Info("测试消息 %d", "ERROR_DEFAULT", index)
			wg.Done()
		}(i)
	}
	wg.Wait()

	t.Log("done")

}
