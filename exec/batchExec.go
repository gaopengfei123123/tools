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

func worker(id int, jobs <-chan int, result chan<- int) {
	for j := range jobs {
		logs.Trace("worker: %d, processing job: %d", id, j)
		time.Sleep(time.Second)
		result <- j * 2
	}
}

func WorkPoolDemo() {
	jobs := make(chan int, 100)
	results := make(chan int, 100)

	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}

	for j := 1; j <= 9; j++ {
		jobs <- j
	}
	close(jobs)
	for a := 1; a <= 9; a++ {
		res := <-results
		logs.Trace("res: %v", res)
	}
}

// CallTask 发起并发执行任务
type CallTask struct {
	MaxTime  time.Duration // 最大执行时间
	TaskList []CallBody    // 执行的任务列表
	Ctx      context.Context
	WorkNum  int
}

// CallBody 单个函数的请求体
type CallBody struct {
	FuncName interface{}   // 要执行的函数本体
	Params   []interface{} // 要传的函数参数, 需要类型和数量都保持一致
	Result   []interface{} // 函数返回的所有结果(包括常用 error), 包括业务返回的错误
	Index    int           // 传入时的顺序
	Err      error         // 函数在 callBack 里调用时的错误
}

// GetResult 将返回结果以反射的方式赋值给传入参数
func (cb *CallBody) GetResult(returnItems ...interface{}) error {
	if cb.Err != nil {
		return cb.Err
	}
	return tools.InterfaceToResult(cb.Result, returnItems...)
}

// Worker 工作线程
func Worker(ctx context.Context, jobs <-chan CallBody, result chan<- CallBody) {
	logs.Info("ctx index: %v", ctx.Value("workIndex"))
LOOP:
	for {
		select {
		case <-ctx.Done():
			logs.Error("work stop, ctx index: %v", ctx.Value("workIndex"))
			break LOOP
		case curJob, ok := <-jobs:
			if !ok {
				logs.Error("work closed: %v", ok)
				break LOOP
			}
			funcRes, funErr := CallFunc(curJob)
			curJob.Err = funErr
			curJob.Result = funcRes
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

func (task *CallTask) AddTask(job CallBody) *CallTask {
	if task.TaskList == nil {
		task.TaskList = make([]CallBody, 0, 10)
	}
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

	// 设置个默认的工作线程
	if task.WorkNum == 0 {
		task.WorkNum = 5
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
		for i := 0; i < max; i++ {
			res, ok := <-resultChan
			if !ok {
				logs.Error("chan is empty")
				return
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
			break LOOP
		}
	}

	return nil
}
