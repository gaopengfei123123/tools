package _739dailyTemperatures

import (
	"github.com/astaxie/beego/logs"
	"testing"
)

func TestCase(t *testing.T) {
	demo := []int{73, 74, 75, 71, 69, 72, 76, 73}
	expect := []int{1, 1, 4, 2, 1, 1, 0, 0}
	logs.Info("demo: %v expect: %v, res: %v", demo, expect, dailyTemperatures(demo))
}

func TestCase2(t *testing.T) {
	demo := []int{30, 40, 50, 60}
	expect := []int{1, 1, 1, 0}
	logs.Info("demo: %v expect: %v, res: %v", demo, expect, dailyTemperatures(demo))
}

func TestCase3(t *testing.T) {
	demo := []int{30, 60, 90}
	expect := []int{1, 1, 0}
	logs.Info("demo: %v expect: %v, res: %v", demo, expect, dailyTemperatures(demo))
}
