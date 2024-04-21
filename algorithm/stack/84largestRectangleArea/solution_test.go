package _84largestRectangleArea

import (
	"github.com/astaxie/beego/logs"
	"testing"
)

func TestCase(t *testing.T) {
	demo := []int{2, 1, 5, 6, 2, 3}
	expect := 10
	logs.Info("demo: %v expect: %v, res: %v", demo, expect, largestRectangleArea(demo))
}

func TestCase2(t *testing.T) {
	demo := []int{2, 4}
	expect := 4
	logs.Info("demo: %v expect: %v, res: %v", demo, expect, largestRectangleArea(demo))
}
