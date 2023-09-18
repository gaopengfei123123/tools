package tools

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"testing"
)

// 将 []interface{} 转成到对应的变量上
func TestInterfaceToResult(t *testing.T) {
	var item1 string
	var item2 error

	result := []interface{}{
		"123123",
		fmt.Errorf("报错%v", 233),
	}

	err := InterfaceToResult(result, &item1, &item2)

	logs.Trace("item1: %#+v, item2: %#+v, err: %v", item1, item2, err)
}
