package _42trap

import (
	"github.com/astaxie/beego/logs"
	"testing"
)

func TestCase(t *testing.T) {
	demo := []int{0, 1, 0, 2, 1, 0, 1, 3, 2, 1, 2, 1}
	expect := 6
	logs.Info("demo: %v expect: %v, res: %v", demo, expect, trap(demo))
}

func TestCase2(t *testing.T) {
	demo := []int{4, 2, 0, 3, 2, 5}
	expect := 9
	logs.Info("demo: %v expect: %v, res: %v", demo, expect, trap(demo))
}
