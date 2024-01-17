package convert

import (
	"github.com/astaxie/beego/logs"
	"testing"
)

func TestDiv(t *testing.T) {
	for i := 1; i < 30; i++ {
		t.Logf("row:%v => %v\n", i, Div(i))
	}
}

type Demo1 struct {
	Key1    string
	Key2    int
	key3    int
	KeyList []string
}

type Demo2 struct {
	Key1    string
	Key2    int
	key3    int
	KeyList []string
}

func TestStructAssign(t *testing.T) {
	src := &Demo1{
		Key1: "aaa", Key2: 2, key3: 3, KeyList: []string{"aaa", "bbb"},
	}

	target := new(Demo2)
	StructAssign(target, src)

	b, _ := JSONEncode(target, true)
	logs.Debug("target: %s", b)
}
