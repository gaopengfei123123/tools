package tools

import (
	"context"
	"github.com/astaxie/beego/logs"
	"testing"
	"time"
)

// 调用示例
func TestDemo(t *testing.T) {
	t.Log("start")
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

	// TODO 现在有个问题, 出现函数超时时不能正确返回剩余成功的结果,
	// eg: 整体限时3s, 函数执行时间也是3s, 就会因为 go coroutines 的协程调度而返回不一致
	batchList := CallTask{
		TaskList: []CallBody{
			func1,
			func2,
			func3,
		},
		Ctx:     context.Background(),
		MaxTime: time.Second * 2,
	}
	res, err := BatchExec(batchList)
	logs.Debug("output result:\nresult: %v\nerr: %s", res, err)

	for i := 0; i < len(res); i++ {
		curResult := res[i]
		if i < 2 {
			// DemoFunc 的返回值
			var msg string
			execErr := curResult.GetResult(&msg)
			logs.Info("exec index: %d, execErr: %v res: %#+v", curResult.Index, execErr, msg)
		} else {
			var res map[string]string
			var err error
			execErr := curResult.GetResult(&res, &err)
			logs.Info("exec index: %d, execErr: %v res: %#+v, %#+v", curResult.Index, execErr, res, err)
		}
	}

	t.Log("end")
}

func DemoFunc(msg string, tt int) string {
	time.Sleep(time.Second * time.Duration(tt))
	logs.Debug("hello demo, sleep %v s", time.Second*time.Duration(tt))
	return "hello " + msg
}

func DemoFunc2(p1, p2, p3 string) (res map[string]string, err error) {
	res = make(map[string]string)
	res["p1"] = p1
	res["p2"] = p2
	res["p3"] = p3
	return
}
