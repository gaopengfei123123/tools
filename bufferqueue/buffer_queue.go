package bufferqueue

import (
	"context"
	"fmt"
	"github.com/astaxie/beego/logs"
	"sync"
	"time"
)

type BufferQueue struct {
	ServerID        string
	job             chan interface{} // 任务队列
	Cap             int              // 队列长度
	MaxTime         time.Duration    // 最大执行间隔
	Ctx             context.Context
	cancelFunc      func()                          // 主动执行关闭
	FlushBufferFunc func(messageList []interface{}) // 清空缓存的回调
	Errors          map[string]error
	execMap         sync.Map // 每个worker 的执行状态, 如果获取到值, 说明正在执行,不再重复执行
}

// NewBufferQueue 这里做简单的初始化, 添加默认值, 在 Start 执行之前都能改
func NewBufferQueue() *BufferQueue {
	bf := new(BufferQueue)
	bf.Cap = 1000
	bf.MaxTime = time.Second * 5
	return bf
}

func (bf *BufferQueue) Start(pCtx context.Context) {
	ctx, cancel := context.WithCancel(pCtx)
	bf.Ctx = ctx
	bf.cancelFunc = cancel

	// 启动的默认值
	if bf.Cap == 0 {
		bf.Cap = 1000
	}

	// 启动的默认值
	if bf.MaxTime == 0 {
		bf.MaxTime = time.Second * 5
	}

	bf.job = make(chan interface{}, bf.Cap)
	go bf.worker(1)
}

func (bf *BufferQueue) Stop() {
	logs.Trace("ID: %v stop", bf.ServerID)
	bf.cancelFunc()
}

func (bf *BufferQueue) AddJob(msg interface{}) error {
	if len(bf.job) == bf.Cap {
		return fmt.Errorf("ID: %v 队列已满, 不能写入", bf.ServerID)
	}
	bf.job <- msg
	return nil
}

// FlushBuffer 清空切片缓存, 如果外部传入了清空逻辑, 则优先执行那个
func (bf *BufferQueue) FlushBuffer(messageList []interface{}, workID int) bool {
	logs.Trace("ID: %v flush buffer, len: %v", bf.ServerID, len(messageList))
	_, ok := bf.execMap.Load(workID)
	if ok {
		logs.Trace("func is running, work_id: %v", workID)
		return false
	}
	// 同一个worker 下, 同时间只能有一个 FlushBuffer 在运行
	bf.execMap.Store(workID, 1)
	defer bf.execMap.Delete(workID)

	if bf.FlushBufferFunc != nil {
		logs.Trace("ID: %v run FlushBufferFunc, work_id: %v", bf.ServerID, workID)
		bf.FlushBufferFunc(messageList)
	}
	logs.Trace("end flush buffer, len: %v", len(messageList))
	return true
}

func (bf *BufferQueue) worker(workID int) {
	logs.Trace("ID: %v  start worker, id: %v", bf.ServerID, workID)
	// 定时任务清理缓存
	ticker := time.NewTicker(bf.MaxTime)
	// 创建一个缓冲的数组, 提供给 FlushBufferFunc 使用
	bufferArr := make([]interface{}, 0, bf.Cap)
	defer ticker.Stop()

LOOP:
	for {
		select {
		case <-ticker.C:
			//logs.Trace("worker:%v is working", workID)
			done := bf.FlushBuffer(bufferArr, workID)
			if done {
				// 外部清空切片
				bufferArr = bufferArr[:0]
			}
		case jb := <-bf.job:
			//logs.Trace("get job")
			bufferArr = append(bufferArr, jb)
			// 获取队列消息, 到达上限后, 就清空一次
			if len(bufferArr) >= bf.Cap {
				logs.Trace("buffer is full, force clean")
				done := bf.FlushBuffer(bufferArr, workID)
				if done {
					// 外部清空切片
					bufferArr = bufferArr[:0]
				}
			}

		case <-bf.Ctx.Done():
			logs.Trace("close channel")
			done := bf.FlushBuffer(bufferArr, workID)
			if done {
				// 外部清空切片
				bufferArr = bufferArr[:0]
			}
			break LOOP
		}
	}

	done := bf.FlushBuffer(bufferArr, workID)
	if done {
		// 外部清空切片
		bufferArr = bufferArr[:0]
	}
	logs.Trace("work is done")
	return
}
