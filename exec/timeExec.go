package exec

import (
	"context"
	"time"
)

// FuncWithTimeout 带超时时间的函数
func FuncWithTimeout(maxTime time.Duration, funcName interface{}, args ...interface{}) *CallBody {
	ctx, _ := context.WithTimeout(context.TODO(), maxTime)
	resChan := make(chan *CallBody, 1)

	task := &CallBody{
		FuncName: funcName,
		Params:   args,
	}

	go func(tsk *CallBody) {
		res, err := CallFunc(*tsk)
		tsk.Result = res
		tsk.Err = err
		tsk.Status = StatusDone // 标记执行完成
		resChan <- tsk
	}(task)

LOOP:
	for {
		select {
		case <-ctx.Done():
			task.Status = StatusOvertime // 标记超时
			break LOOP
		case tsk, _ := <-resChan:
			task = tsk
			break LOOP
		}
	}
	return task
}
