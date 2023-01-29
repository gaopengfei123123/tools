package bufferqueue

import (
	"context"
	"fmt"
	"github.com/astaxie/beego/logs"
	"sync"
	"time"
)

type Logger interface {
	Trace(format string, v ...interface{})
}

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
	debug           bool
	Logger          Logger // 日志
}

// NewBufferQueue 这里做简单的初始化, 添加默认值, 在 Start 执行之前都能改
func NewBufferQueue() *BufferQueue {
	bf := new(BufferQueue)
	bf.Ctx = context.TODO()
	bf.Cap = 1000
	bf.MaxTime = time.Second * 5
	bf.Logger = logs.GetBeeLogger()
	return bf
}

func (bf *BufferQueue) DebugMode() {
	bf.debug = true
}

func (bf *BufferQueue) Log(format string, v ...interface{}) {
	if bf.debug && bf.Logger != nil {
		bf.Logger.Trace(format, v)
	}
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
	bf.Log("ID: %v stop", bf.ServerID)
	bf.cancelFunc()
}

func (bf *BufferQueue) AddJob(msg interface{}) error {
	var err error

	if bf.job == nil {
		return fmt.Errorf("please exec BufferQueue.")
	}

	logs.Info("job channel: %v", bf.job)
	bf.job <- msg
	return err
}

// FlushBuffer 清空切片缓存, 如果外部传入了清空逻辑, 则优先执行那个
func (bf *BufferQueue) FlushBuffer(messageList []interface{}, workID int) bool {
	bf.Log("ID: %v flush buffer, len: %v", bf.ServerID, len(messageList))
	//_, ok := bf.execMap.Load(workID)
	//
	//if ok {
	//	bf.Log("func is running, work_id: %v", workID)
	//	return false
	//}
	//// 同一个worker 下, 同时间只能有一个 FlushBuffer 在运行
	//bf.execMap.Store(workID, 1)
	//defer bf.execMap.Delete(workID)

	if bf.FlushBufferFunc != nil {
		bf.Log("ID: %v run FlushBufferFunc, work_id: %v", bf.ServerID, workID)
		bf.FlushBufferFunc(messageList)
	}
	bf.Log("end flush buffer, len: %v", len(messageList))
	return true
}

func (bf *BufferQueue) worker(workID int) {
	bf.Log("ID: %v  start worker, id: %v", bf.ServerID, workID)
	// 定时任务清理缓存
	ticker := time.NewTicker(bf.MaxTime)
	// 创建一个缓冲的数组, 提供给 FlushBufferFunc 使用
	bufferArr := make([]interface{}, 0, bf.Cap)
	defer ticker.Stop()

LOOP:
	for {
		select {
		case <-ticker.C:
			//bf.Log("worker:%v is working", workID)
			done := bf.FlushBuffer(bufferArr, workID)
			if done {
				// 外部清空切片
				bufferArr = bufferArr[:0]
			}
		case jb := <-bf.job:
			//bf.Log("get job")
			bufferArr = append(bufferArr, jb)
			// 获取队列消息, 到达上限后, 就清空一次
			if len(bufferArr) >= bf.Cap {
				bf.Log("buffer is full, force clean")
				done := bf.FlushBuffer(bufferArr, workID)
				if done {
					// 外部清空切片
					bufferArr = bufferArr[:0]
				}
			}

		case <-bf.Ctx.Done():
			bf.Log("close channel")
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
	bf.Log("work is done")
	return
}
