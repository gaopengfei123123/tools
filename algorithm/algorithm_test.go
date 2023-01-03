package algorithm

import (
	"github.com/astaxie/beego/logs"
	"testing"
)

func TestHammingWeight(t *testing.T) {
	logs.Info("11:=", hammingWeight(11))
}

func TestCheckNuggets(t *testing.T) {
	logs.Info("7 := ", CheckNuggets(7))
	logs.Info("25 := ", CheckNuggets(25))
	logs.Info("29 := ", CheckNuggets(29))
}

func TestFindKthLargest(t *testing.T) {
	//logs.Info("[3,2,1,5,6,4], k = 2, res:", FindKthLargest([]int{3, 2, 1, 5, 6, 4}, 2))
	logs.Info("[3,2,1,5,6,4], k = 2, res:", _quickSort([]int{3, 2, 1, 5, 6, 4}, 0, 5))
}

func TestFindKthLargest2(t *testing.T) {
	v := float64(1005328)
	logs.Info("value: %v", v)
	logs.Info("value: %0.f", v)

}
