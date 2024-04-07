package midianFinder

import "container/heap"

// 原文地址 https://leetcode.cn/problems/find-median-from-data-stream/solutions/2361972/295-shu-ju-liu-de-zhong-wei-shu-dui-qing-gmdo/?envType=study-plan-v2&envId=top-100-liked

/**
解题思路:
1. 维护左右两个堆, 左边为最小堆A, 右边为最大堆B
2. 插入元素时, 当数量 A=B 时, 新元素和 A 最小值 中取两者最小值, 插入到 B 中, 当数量 A!=B 时, 新元素和B 最大值取最大值, 插入到A 中
3. 获取中位数时, 如果数量 A=B 时, 中位数是 (A最小值+B 最大值)/2, 如果 A!=B, 则中位数是 B 最大值
*/

type MedianFinder struct {
	LeftHp  *MaxHeap // 左边
	RightHp *MinHeap // 右边
	// 整体数组左大于右
}

func Constructor() MedianFinder {
	return MedianFinder{&MaxHeap{}, &MinHeap{}}
}

func (this *MedianFinder) AddNum(num int) {
	// 确保新加入的值始终是 left < num < right
	// 优先加入left, 其次right
	if num < this.LeftHp.Head() && this.LeftHp.Len() != 0 {
		x := heap.Pop(this.LeftHp)
		heap.Push(this.LeftHp, num)
		num = x.(int)
	}

	if num > this.RightHp.Head() && this.RightHp.Len() != 0 {
		x := heap.Pop(this.RightHp)
		heap.Push(this.RightHp, num)
		num = x.(int)
	}

	if this.LeftHp.Len() == this.RightHp.Len() {
		heap.Push(this.LeftHp, num)
	} else {
		heap.Push(this.RightHp, num)
	}
}

//  1
// 1 | 2
// 1 2 | 3

func (this *MedianFinder) FindMedian() float64 {
	if this.LeftHp.Len() == this.RightHp.Len() {
		return (float64(this.RightHp.Head()) + float64(this.LeftHp.Head())) / 2
	}
	return float64(this.LeftHp.Head())
}

/**
 * Your MedianFinder object will be instantiated and called as such:
 * obj := Constructor();
 * obj.AddNum(num);
 * param_2 := obj.FindMedian();
 */

// MinHeap 最小堆
type MinHeap []int

func (hp MinHeap) Len() int           { return len(hp) }
func (hp MinHeap) Less(i, j int) bool { return hp[i] < hp[j] }
func (hp MinHeap) Swap(i, j int)      { hp[i], hp[j] = hp[j], hp[i] }
func (hp *MinHeap) Push(x interface{}) {
	*hp = append(*hp, x.(int))
}

func (hp *MinHeap) Pop() interface{} {
	n := len(*hp)
	x := (*hp)[n-1]
	*hp = (*hp)[:n-1]
	return x
}
func (hp MinHeap) Head() int {
	if len(hp) == 0 {
		return 0
	}
	return hp[0]
}

// MaxHeap 最大堆
type MaxHeap []int

func (hp MaxHeap) Len() int           { return len(hp) }
func (hp MaxHeap) Less(i, j int) bool { return hp[i] > hp[j] }
func (hp MaxHeap) Swap(i, j int)      { hp[i], hp[j] = hp[j], hp[i] }
func (hp *MaxHeap) Push(x interface{}) {
	*hp = append(*hp, x.(int))
}

func (hp *MaxHeap) Pop() interface{} {
	n := len(*hp)
	x := (*hp)[n-1]
	*hp = (*hp)[:n-1]
	return x
}

func (hp MaxHeap) Head() int {
	if len(hp) == 0 {
		return 0
	}
	return hp[0]
}
