package timeutil

import (
	"github.com/astaxie/beego/logs"
	"testing"
	"time"
)

func TestTimeCnt(t *testing.T) {
	fn := TimeCnt("testTime")
	time.Sleep(time.Second)

	_, _, msg := fn()
	logs.Info(msg)
}
