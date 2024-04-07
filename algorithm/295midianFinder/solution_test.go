package _95midianFinder

import (
	"container/heap"
	"github.com/astaxie/beego/logs"
	"testing"
)

func TestConstructor(t *testing.T) {
	maxhp := &MaxHeap{}

	heap.Push(maxhp, 5)
	heap.Push(maxhp, 2)
	heap.Push(maxhp, 1)
	heap.Push(maxhp, 4)
	x := heap.Pop(maxhp)
	logs.Info("pop: %v", x)
	logs.Info("maxHp: %v", maxhp)

	minHp := &MinHeap{}
	heap.Push(minHp, 5)
	heap.Push(minHp, 2)
	heap.Push(minHp, 1)
	heap.Push(minHp, 4)
	x = heap.Pop(minHp)
	logs.Info("pop: %v", x)
	logs.Info("minHp: %v", minHp)
}

func TestConstructor2(t *testing.T) {
	body := Constructor()

	body.AddNum(1)
	logs.Info("body: %v %v", body.LeftHp, body.RightHp)
	body.AddNum(2)
	logs.Info("body: %v %v", body.LeftHp, body.RightHp)
	body.AddNum(3)
	logs.Info("median: %v", body.FindMedian())
	body.AddNum(4)
	logs.Info("median: %v", body.FindMedian())
	logs.Info("body: %v %v", body.LeftHp, body.RightHp)
}
