package tools

import (
	"context"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"reflect"
	"time"
)

/**
并发执行函数, 使用示例见 combineExec_test.go 中的测试用例

TODO 升级版的方法请异步 asyncExec_test.go

目前最佳试用场景是在 for 循环中重复执行的函数可以使用, 或者不需要接受返回结果的

其他场景解析不同的函数返回数据会很繁琐

目前使用时碰到的问题:
1. 如果超时设置为 3s, func 耗时也是 3s, 那因为 goroutine 调度机制下, 将随机返回数据, 这个边界情况需要考虑
2. 参数传的值中包含指针的时候, 需要留意批量执行时的传参是否符合预期, 大概率需要拷贝出一个新的参数传进去
*/

// CallTask 发起并发执行任务
type CallTask struct {
	MaxTime  time.Duration
	TaskList []CallBody
	Ctx      context.Context
}

// CallBody 单个函数的请求体
type CallBody struct {
	FuncName interface{}   // 要执行的函数本体
	Params   []interface{} // 要传的函数参数, 需要类型和数量都保持一致
}

type ResponseBody struct {
	Result []interface{} // 函数返回的所有结果(包括常用 error), 包括业务返回的错误
	Index  int           // 传入时的顺序
	Err    error         // 函数在 callBack 里调用时的错误
}

// GetResult 将返回结果以反射的方式赋值给传入参数
func (resp *ResponseBody) GetResult(returnItems ...interface{}) error {
	if resp.Err != nil {
		return resp.Err
	}
	return InterfaceToResult(resp.Result, returnItems...)
}

// BatchExec 并发执行, 需限制总的并发数量
func BatchExec(task CallTask) (resultList []ResponseBody, err error) {
	if len(task.TaskList) == 0 {
		return
	}

	if len(task.TaskList) > 500 {
		err = fmt.Errorf("go coroutines limit 500")
		return
	}

	// 设置个默认超时时间
	if task.MaxTime == 0 {
		task.MaxTime = time.Second * 10
	}

	// 创建阻塞通道, 为了计算超时时间
	resultChan := make(chan ResponseBody, len(task.TaskList))
	timeoutCtx, cancel := context.WithTimeout(task.Ctx, task.MaxTime)
	defer close(resultChan)
	defer cancel()

	for i := 0; i < len(task.TaskList); i++ {
		logs.Debug("params index %d : %#+v \n", i, task.TaskList[i].Params)
		//result, err := CallFunc(task.TaskList[i])
		// 这里开协程进行批量调用
		go func(index int, task CallBody, res chan<- ResponseBody) {
			result, resErr := CallFunc(task)
			// 将结果和错误信息返回
			res <- ResponseBody{
				Result: result,
				Err:    resErr,
				Index:  index,
			}
			return
		}(i, task.TaskList[i], resultChan)
	}

	// 用来标记已经读取全部返回结果
	done := make(chan struct{}, 1)
	defer close(done)
	result := make([]ResponseBody, 0, len(task.TaskList))
	// 并行读取 这里用的闭包的外部变量, 在时间截止前能收到多少消息就存多少
	go func(max int, ctx context.Context) {
		for i := 0; i < max; i++ {
			res, ok := <-resultChan
			if !ok {
				logs.Error("chan closed")
				return
			}
			logs.Debug("get func result: %#+v", res)
			result = append(result, res)
		}
		done <- struct{}{}
	}(len(task.TaskList), timeoutCtx)

LOOP:
	for {
		select {
		case <-done:
			logs.Debug("exec end")
			break LOOP
		case <-timeoutCtx.Done():
			logs.Debug("timeout")
			break LOOP
		}
	}
	logs.Debug("final result: %#+v", result)
	resultList = make([]ResponseBody, len(task.TaskList))
	for i := range resultList {
		resultList[i].Err = errors.New("func exec timeout")
		resultList[i].Index = i
	}

	for _, v := range result {
		resultList[v.Index] = v
	}
	return resultList, nil
}

// SyncBatchExec 同步执行, 用于不方便并发执行的业务
func SyncBatchExec(task CallTask) (resultList []ResponseBody, err error) {
	if len(task.TaskList) == 0 {
		return
	}

	// 设置个默认超时时间, 串行执行
	if task.MaxTime == 0 {
		task.MaxTime = time.Second * 600
	}

	// 创建阻塞通道, 为了计算超时时间
	resultChan := make(chan ResponseBody, len(task.TaskList))
	timeoutCtx, cancel := context.WithTimeout(task.Ctx, task.MaxTime)
	defer close(resultChan)
	defer cancel()

	for i := 0; i < len(task.TaskList); i++ {
		logs.Debug("params index %d : %#+v \n", i, task.TaskList[i].Params)
		//result, err := CallFunc(task.TaskList[i])
		// 这里串行批量调用
		func(index int, task CallBody, res chan<- ResponseBody) {
			result, resErr := CallFunc(task)
			// 将结果和错误信息返回
			res <- ResponseBody{
				Result: result,
				Err:    resErr,
				Index:  index,
			}
			return
		}(i, task.TaskList[i], resultChan)
	}

	// 用来标记已经读取全部返回结果
	done := make(chan struct{}, 1)
	defer close(done)
	result := make([]ResponseBody, 0, len(task.TaskList))
	// 并行读取 这里用的闭包的外部变量, 在时间截止前能收到多少消息就存多少
	go func(max int, ctx context.Context) {
		defer func() {
			if panicErr := recover(); panicErr != nil {
				err = fmt.Errorf("%v", panicErr)
			}
		}()

		for i := 0; i < max; i++ {
			res, ok := <-resultChan
			if !ok {
				logs.Error("chan closed")
				return
			}
			logs.Debug("get func result: %#+v", res)
			result = append(result, res)
		}
		done <- struct{}{}
	}(len(task.TaskList), timeoutCtx)

LOOP:
	for {
		select {
		case <-done:
			logs.Debug("exec end")
			break LOOP
		case <-timeoutCtx.Done():
			logs.Debug("timeout")
			break LOOP
		}
	}
	logs.Debug("final result: %#+v", result)
	resultList = make([]ResponseBody, len(task.TaskList))
	for i := range resultList {
		resultList[i].Err = errors.New("func exec timeout")
		resultList[i].Index = i
	}

	for _, v := range result {
		resultList[v.Index] = v
	}
	return resultList, nil
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
