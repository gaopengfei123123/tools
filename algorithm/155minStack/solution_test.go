package _155minStack

import (
	"github.com/astaxie/beego/logs"
	"testing"
)

func TestConstructor(t *testing.T) {
	stack := Constructor()
	stack.Push(-2)
	stack.Push(0)
	stack.Push(-3)
	stack.Push(1)

	//stack.Pop()

	logs.Info("current Stack: %v, top: %v, min: %v", stack, stack.Top(), stack.GetMin())
	stack.Pop()
	logs.Info("current Stack: %v, top: %v, min: %v", stack, stack.Top(), stack.GetMin())
	stack.Pop()
	logs.Info("current Stack: %v, top: %v, min: %v", stack, stack.Top(), stack.GetMin())
}
