package bufferqueue

import (
	"context"
	"fmt"
	"github.com/astaxie/beego/logs"
	"testing"
	"time"
)

// 这里对程序进行简单的测试
func TestBufferQueue_Start(t *testing.T) {
	buffer := NewBufferQueue()
	timeout := time.After(time.Second * 30)
	ticker := time.NewTicker(time.Second * 1)

	ctx, cancel := context.WithCancel(context.TODO())
	buffer.Start(ctx)

	i := 0
LOOP:
	for {
		select {
		case <-ticker.C:
			err := buffer.AddJob(fmt.Sprintf("a new Job, index: %d", i))
			i += 1
			if err != nil {
				logs.Trace("addJob err: %v", err)
			}

			//// 模拟主动退出的情况
			//if i == 5 {
			//	cancel()
			//}
		case <-timeout:
			logs.Trace("timeout, stop")
			break LOOP
		}
	}
	cancel()
	time.Sleep(time.Second * 1)
	logs.Trace("finish")
}

// 这里注入逻辑函数, 看一下执行效果
func TestBufferQueue_FlushBuffer(t *testing.T) {
	buffer := NewBufferQueue()
	timeout := time.After(time.Second * 15)
	ticker := time.NewTicker(time.Second * 3)

	ctx, cancel := context.WithCancel(context.TODO())

	// 这里注入了外部执行的方法
	buffer.FlushBufferFunc = demoFunction

	buffer.Start(ctx)
	i := 0
LOOP:
	for {
		select {
		case <-ticker.C:
			err := buffer.AddJob(fmt.Sprintf("a new Job, index: %d", i))
			i += 1
			if err != nil {
				logs.Trace("addJob err: %v", err)
			}
		case <-timeout:
			logs.Trace("timeout, stop")
			break LOOP
		}
	}
	cancel()
	time.Sleep(time.Second * 1)
	logs.Trace("finish")
}

func demoFunction(message []interface{}) {
	logs.Trace("demoFunc: %v", message)
}

func TestBufferQueue_AddJob(t *testing.T) {
	buffer := NewBufferQueue()
	buffer.Start(context.TODO())

	err := buffer.AddJob("xxx")
	logs.Trace("err: %v", err)
}

// 模拟主动退出
func TestBufferQueue_Stop(t *testing.T) {
	buffer := NewBufferQueue()
	timeout := time.After(time.Second * 30)
	ticker := time.NewTicker(time.Second * 1)
	ctx, cancel := context.WithCancel(context.TODO())
	buffer.Start(ctx)

	i := 0
LOOP:
	for {
		select {
		case <-ticker.C:
			err := buffer.AddJob(fmt.Sprintf("a new Job, index: %d", i))
			i += 1
			if err != nil {
				logs.Trace("addJob err: %v", err)
				break LOOP
			}

			// 模拟主动退出的情况
			if i == 5 {
				buffer.Stop()
			}
		case <-timeout:
			logs.Trace("timeout, stop")
			break LOOP
		}
	}
	cancel()
	time.Sleep(time.Second * 1)
	logs.Trace("finish")
}
