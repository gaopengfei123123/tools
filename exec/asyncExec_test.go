package exec

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"testing"
	"time"
)

// 单协程执行示例
func TestCallTask_BatchExec(t *testing.T) {
	task := &CallTask{}
	task.WorkNum = 1 // 工作协程为1, 相当于同步执行
	for i := 0; i < 10; i++ {
		func1 := CallBody{
			FuncName: DemoFunc,
			Params: []interface{}{
				fmt.Sprintf("index: %v", i),
				1,
			},
		}
		task.AddTask(func1)
	}

	logs.Trace("task: %#+v", task)
	task.BatchExec()
	logs.Trace("task: %#+v", task)
}

// 多协程执行示例
func TestCallTask_BatchExec2(t *testing.T) {
	task := &CallTask{}
	task.WorkNum = 3
	for i := 0; i < 10; i++ {
		func1 := CallBody{
			FuncName: DemoFunc,
			Params: []interface{}{
				fmt.Sprintf("index: %v", i),
				1,
			},
		}
		task.AddTask(func1)
	}

	logs.Trace("task: %#+v", task)
	task.BatchExec()
	logs.Trace("task: %#+v", task)

	for i := range task.TaskList {
		var msg string
		err := task.TaskList[i].GetResult(&msg)
		logs.Trace("func result => index: %d, msg: %s, FuncErr: %v", i, msg, err)
	}
}

// 多函数类型执行示例, 以及获取返回值示例
func TestCallTask_BatchExec3(t *testing.T) {
	func1 := CallBody{
		FuncName: DemoFunc,
		Params: []interface{}{
			"word", 3,
		},
	}

	func2 := CallBody{
		FuncName: DemoFunc,
		Params: []interface{}{
			"GPF", 1,
		},
	}

	func3 := CallBody{
		FuncName: DemoFunc2,
		Params: []interface{}{
			"XXX", "YYY", "ZZZ",
		},
	}

	task := &CallTask{}
	task.MaxTime = time.Second * 2 // 设置最大超时时间 2s
	task.AddTask(func1).AddTask(func2).AddTask(func3)

	funcErr := task.BatchExec()

	if funcErr != nil {
		logs.Error("batchExecErr: %v", funcErr)
		return
	}

	for i := 0; i < len(task.TaskList); i++ {
		curResult := task.TaskList[i]
		if i < 2 {
			// DemoFunc 的返回值
			var msg string
			execErr := curResult.GetResult(&msg)
			logs.Trace("exec index: %d, execErr: %v res: %#+v", curResult.Index, execErr, msg)
		} else {
			var res map[string]string
			var err error
			execErr := curResult.GetResult(&res, &err)
			logs.Trace("exec index: %d, execErr: %v res: %#+v, %#+v", curResult.Index, execErr, res, err)
		}
	}

}

func DemoFunc(msg string, tt int) string {
	time.Sleep(time.Second * time.Duration(tt))
	logs.Debug("DemoFunc: hello %s, sleep %v s", msg, time.Second*time.Duration(tt))
	return "hello " + msg
}

func DemoFunc2(p1, p2, p3 string) (res map[string]string, err error) {
	res = make(map[string]string)
	res["p1"] = p1
	res["p2"] = p2
	res["p3"] = p3
	return
}

func TestFuncWithTimeout(t *testing.T) {
	var msg string

	err := FuncWithTimeout(time.Second*2, DemoFunc, "hello world", 1).GetResult(&msg)
	logs.Info("err: %v, result: %v", err, msg)

	var msg2 string
	err = FuncWithTimeout(time.Second*2, DemoFunc, "hello world", 3).GetResult(&msg2)
	logs.Info("err: %v, result: %v", err, msg2)
}
