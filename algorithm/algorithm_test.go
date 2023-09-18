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

func TestFindKthLargest(t *testing.T) {
	//logs.Trace("[3,2,1,5,6,4], k = 2, res:", FindKthLargest([]int{3, 2, 1, 5, 6, 4}, 2))
	logs.Trace("[3,2,1,5,6,4], k = 2, res:", _quickSort([]int{3, 2, 1, 5, 6, 4}, 0, 5))
}

func TestFindKthLargest2(t *testing.T) {
	v := float64(1005328)
	logs.Trace("value: %v", v)
	logs.Trace("value: %0.f", v)

}
