package algorithm

import (
	"github.com/astaxie/beego/logs"
	"testing"
)

func TestHammingWeight(t *testing.T) {
	logs.Trace("11:=", hammingWeight(11))
}

func TestCheckNuggets(t *testing.T) {
	logs.Trace("7 := ", CheckNuggets(7))
	logs.Trace("25 := ", CheckNuggets(25))
	logs.Trace("29 := ", CheckNuggets(29))
}

func TestClimbStairs(t *testing.T) {
	logs.Info("3 := %v", climbStairs(3))
	logs.Info("6 := %v", climbStairs(6))
}

//func TestFindKthLargest(t *testing.T) {
//	//logs.Trace("[3,2,1,5,6,4], k = 2, res:", FindKthLargest([]int{3, 2, 1, 5, 6, 4}, 2))
//	logs.Trace("[3,2,1,5,6,4], k = 2, res:", _quickSort([]int{3, 2, 1, 5, 6, 4}, 0, 5))
//}

func TestFindKthLargest2(t *testing.T) {
	v := float64(1005328)
	logs.Trace("value: %v", v)
	logs.Trace("value: %0.f", v)

}

func TestFindMaxConsecutiveOnes(t *testing.T) {
	demo := []int{1, 1, 0, 1, 1, 1}
	expect := 3
	logs.Info("demo: %v expect: %v, res: %v", demo, expect, FindMaxConsecutiveOnes(demo))

	demo = []int{1, 0, 1, 1, 0, 1}
	expect = 2
	logs.Info("demo: %v expect: %v, res: %v", demo, expect, FindMaxConsecutiveOnes(demo))
}

func TestFindPoisonedDuration(t *testing.T) {
	demo := []int{1, 4}
	duration := 2
	expect := 4
	logs.Info("demo: %v, duration: %v, expect: %v, res: %v", demo, duration, expect, FindPoisonedDuration(demo, duration))

	demo = []int{1, 2, 3, 4}
	duration = 2
	expect = 5
	logs.Info("demo: %v, duration: %v, expect: %v, res: %v", demo, duration, expect, FindPoisonedDuration(demo, duration))
}

func TestFindKthLargest(t *testing.T) {
	demo := []int{3, 2, 1, 5, 6, 4}
	kth := 2
	expect := 5
	logs.Info("demo: %v, kth: %v, expect: %v, res: %v", demo, kth, expect, FindKthLargest(demo, kth))

	demo = []int{3, 2, 1, 5, 6, 4}
	kth = 3
	expect = 4
	logs.Info("demo: %v, kth: %v, expect: %v, res: %v", demo, kth, expect, FindKthLargest(demo, kth))

}

func TestBuildHeap(t *testing.T) {
	demo := []int{3, 2, 1, 5, 6, 4}
	heapSize := len(demo)
	logs.Info("before: %v", demo)

	buildMaxHeap(demo, heapSize)
	logs.Info("after: %v size: %v", demo, heapSize)

	heapSize--
	demo[0], demo[len(demo)-1] = demo[len(demo)-1], demo[0]
	maxHeapify(demo, 0, heapSize)
	logs.Info("after del 1: %v, size: %v", demo, heapSize)

	heapSize--
	demo[0], demo[len(demo)-2] = demo[len(demo)-2], demo[0]
	maxHeapify(demo, 0, heapSize)
	logs.Info("after del 1: %v, size: %v", demo, heapSize)
}
