package topKFrequent

import (
	"github.com/astaxie/beego/logs"
	"testing"
)

func TestTopKFrequent(t *testing.T) {
	demo := []int{1, 1, 1, 2, 2, 3}
	kth := 2
	expect := []int{1, 2}
	logs.Info("demo: %v, kth: %v, expect: %v, res: %v", demo, kth, expect, TopKFrequent(demo, kth))
}

func TestTopKFrequent2(t *testing.T) {
	demo := []int{1}
	kth := 1
	expect := []int{1}
	logs.Info("demo: %v, kth: %v, expect: %v, res: %v", demo, kth, expect, TopKFrequent(demo, kth))
}

func TestTopKFrequent3(t *testing.T) {
	demo := []int{-1, -1}
	kth := 1
	expect := []int{-1}
	logs.Info("demo: %v, kth: %v, expect: %v, res: %v", demo, kth, expect, TopKFrequent(demo, kth))
}

func TestTopKFrequent4(t *testing.T) {
	demo := []int{1, 2}
	kth := 2
	expect := []int{1, 2}
	logs.Info("demo: %v, kth: %v, expect: %v, res: %v", demo, kth, expect, TopKFrequent(demo, kth))

}
