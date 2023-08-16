package exec

import (
	"context"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/gaopengfei123123/tools"
	"github.com/pkg/errors"
	"reflect"
	"time"
)

// 关于线程池简单示例: https://learnku.com/docs/gobyexample/2020/work-pool-worker/6285

/**
并发执行函数, 使用示例见 asyncExec_test.go 中的测试用例

目前最佳试用场景是在 for 循环中重复执行的函数可以使用, 或者不需要接受返回结果的

其他场景解析不同的函数返回数据会很繁琐

新增功能:
1. 添加线程控制, 可以通过 CallTask.WorkNum 控制并发数量
2. 优化

目前使用时碰到的问题:
1. 如果超时设置为 3s, func 耗时也是 3s, 那因为 goroutine 调度机制下, 将随机返回数据, 这个边界情况需要考虑
2. 参数传的值中包含指针的时候, 需要留意批量执行时的传参是否符合预期, 大概率需要拷贝出一个新的参数传进去
*/

//函数执行状态

const StatusDone = 1
const StatusWait = 0
const StatusOvertime = 2

// CallTask 发起并发执行任务
type CallTask struct {
	MaxTime  time.Duration // 最大执行时间
	TaskList []CallBody    // 执行的任务列表
	Ctx      context.Context
	WorkNum  int
}

// CallBody 单个函数的请求体
type CallBody struct {
	Status   int           // 执行状态
	Index    int           // 传入时的顺序
	FuncName interface{}   // 要执行的函数本体
	Params   []interface{} // 要传的函数参数, 需要类型和数量都保持一致
	Result   []interface{} // 函数返回的所有结果(对应多参数返回)
	Err      error         // 函数在 callBack 里调用时的错误
}

// GetResult 将返回结果以反射的方式赋值给传入参数
func (cb *CallBody) GetResult(returnItems ...interface{}) error {
	if cb.Err != nil {
		return cb.Err
	}

	if cb.Status == StatusWait {
		return fmt.Errorf("func not exec yet, index: %v", cb.Index)
	}

	if cb.Status == StatusOvertime {
		return fmt.Errorf("func exec overtime, index: %v", cb.Index)
	}

	return tools.InterfaceToResult(cb.Result, returnItems...)
}

// Worker 工作线程
func Worker(ctx context.Context, jobs <-chan CallBody, result chan<- CallBody) {
	logs.Info("start worker, ctx index: %v", ctx.Value("workIndex"))
LOOP:
	for {
		select {
		case <-ctx.Done():
			logs.Trace("worker timeout stop, ctx index: %v", ctx.Value("workIndex"))
			break LOOP
		case curJob, ok := <-jobs:
			if !ok {
				logs.Trace("worker closed, ctx index: %v", ctx.Value("workIndex"))
				break LOOP
			}
			funcRes, funErr := CallFunc(curJob)
			curJob.Err = funErr
			curJob.Result = funcRes
			curJob.Status = StatusDone
			result <- curJob
		}
	}

	return
}

// CallFunc 利用反射动态执行函数
func CallFunc(body CallBody) (result []interface{}, err error) {
	// 校验是否是函数
	if reflect.TypeOf(body.FuncName).Kind() != reflect.Func {
		err = errors.New(fmt.Sprintf("this is not a  func name abort. FuncName: %v", body.FuncName))
		return
	}
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = fmt.Errorf("%v", panicErr)
		}
	}()

	// 执行方法
	f := reflect.ValueOf(body.FuncName)
	// 校验传参值数量
	if len(body.Params) != f.Type().NumIn() {
		err = errors.New("The number of params is not adapted.")
		return
	}
	in := make([]reflect.Value, len(body.Params))
	for k, param := range body.Params {
		in[k] = reflect.ValueOf(param)
	}
	res := f.Call(in)
	result = make([]interface{}, len(res))
	for k, v := range res {
		result[k] = v.Interface()
	}
	return
}

// AddTask 添加工作任务
func (task *CallTask) AddTask(job CallBody) *CallTask {
	if task.TaskList == nil {
		task.TaskList = make([]CallBody, 0, 10)
	}
	curIndex := len(task.TaskList)
	job.Status = StatusWait
	job.Index = curIndex

	task.TaskList = append(task.TaskList, job)
	return task
}

func (task *CallTask) BatchExec() error {
	if len(task.TaskList) == 0 {
		return nil
	}

	if len(task.TaskList) > 5000 {
		return fmt.Errorf("go coroutines limit 5000")
	}

	if task.Ctx == nil {
		task.Ctx = context.Background()
	}

	// 设置个默认超时时间
	if task.MaxTime == 0 {
		task.MaxTime = time.Second * 10
	}

	// 设置个默认的工作线程数量
	if task.WorkNum == 0 {
		task.WorkNum = 5
	}

	// 任务数少于线程数, 则按任务数来
	if task.WorkNum > len(task.TaskList) {
		task.WorkNum = len(task.TaskList)
	}

	jobChan := make(chan CallBody, len(task.TaskList))
	resultChan := make(chan CallBody, len(task.TaskList))
	defer close(resultChan)

	timeoutCtx, cancel := context.WithTimeout(task.Ctx, task.MaxTime)
	defer cancel()

	// 启动工作进程
	for w := 1; w <= task.WorkNum; w++ {
		childCtx := context.WithValue(timeoutCtx, "workIndex", w)
		go Worker(childCtx, jobChan, resultChan)
	}

	// 将任务放入工作队列
	for j := 0; j < len(task.TaskList); j++ {
		cur := task.TaskList[j]
		cur.Index = j
		jobChan <- cur
	}
	close(jobChan)

	// 用来标记已经读取全部返回结果
	done := make(chan struct{}, 1)
	defer close(done)

	// 并行读取结果
	go func(max int) {
		// 如果执行超时, 则需要捕获 send on closed channel 的异常
		defer func() {
			if panicErr := recover(); panicErr != nil {
				logs.Error("panicErr: %s", panicErr)
				return
			}
		}()

		for i := 0; i < max; i++ {
			res, ok := <-resultChan
			if !ok {
				logs.Trace("chan is empty")
				break
			}
			// 将结果返回给原位置
			task.TaskList[res.Index] = res
		}
		done <- struct{}{}
	}(len(task.TaskList))

LOOP:
	for {
		select {
		case <-done:
			logs.Trace("finish")
			break LOOP
		case <-timeoutCtx.Done():
			logs.Trace("timeout")
			// 所有未执行/未返回结果的函数, 按超时处理
			for i := range task.TaskList {
				if task.TaskList[i].Status == StatusWait {
					task.TaskList[i].Status = StatusOvertime
				}
			}
			break LOOP
		}
	}

	return nil
}
