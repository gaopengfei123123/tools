package algorithm

// 2进制中1的个数 https://leetcode.cn/problems/er-jin-zhi-zhong-1de-ge-shu-lcof/
func hammingWeight(num uint32) int {
	var ret int
	for i := 0; i < 32; i++ {
		if 1<<i&num > 0 {
			ret++
		}
	}
	return ret
}

// 爬楼梯问题 https://leetcode.cn/problems/climbing-stairs/submissions/
// 斐波那契类的问题
func climbStairs(n int) int {
	ln1 := 1
	ln2 := 2

	if n <= 1 {
		return ln1
	}

	if n <= 2 {
		return ln2
	}

	for i := 3; i <= n; i++ {
		tmp := ln1 + ln2
		ln1 = ln2
		ln2 = tmp
	}
	return ln2
}

// 1 2 3 5 8 13 21 34

/**
 CheckNuggets 麦乐鸡块问题
问题：在一个平行宇宙中，麦当劳的麦乐鸡块分为 7 块装、13 块装和 29 块装。有一天，你的老板让你出去购买正好为 n 块(0 < n <= 10000)的麦乐鸡块回来，请提供一个算法判断是否可行。
*/

func CheckNuggets(n int) bool {
	tmpList := make(map[int]bool)
	tmpList[7] = true
	tmpList[13] = true
	tmpList[29] = true

	for i := 7; i <= n; i++ {
		if tmpList[i] {
			tmpList[i+7] = true
			tmpList[i+13] = true
			tmpList[i+29] = true
		}
	}

	return tmpList[n]
}

func _quickSort(arr []int, left, right int) []int {
	if len(arr) == 0 {
		return []int{}
	}
	if left < right {
		partitionIndex := partition(arr, left, right)
		_quickSort(arr, left, partitionIndex-1)
		_quickSort(arr, partitionIndex+1, right)
	}
	return arr
}

func partition(arr []int, left, right int) int {
	pivot := left
	index := pivot + 1
	for i := index; i <= right; i++ {
		if arr[i] < arr[pivot] {
			swap(arr, i, index)
			index += 1
		}
	}
	swap(arr, pivot, index-1)
	return index - 1
}

func swap(arr []int, i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}

// ## 数组相关问题, 按题目顺序练习

// FindMaxConsecutiveOnes  1 最大连续1的个数  力扣编号485
func FindMaxConsecutiveOnes(nums []int) int {
	if len(nums) == 0 {
		return 0
	}

	res := 0

	tmp := 0
	for i := 0; i < len(nums); i++ {
		if nums[i] == 1 {
			tmp++
			if tmp > res {
				res = tmp
			}
		}

		if len(nums) <= i+1 {
			return res
		}

		if nums[i+1] != 1 {
			tmp = 0
		}
	}
	return res
}

// FindPoisonedDuration 2 提莫攻击 力扣编号495
func FindPoisonedDuration(timeSeries []int, duration int) int {
	res := 0
	for i, current := range timeSeries {
		if i+1 >= len(timeSeries) {
			res += duration
			break
		}

		next := timeSeries[i+1]
		limit := next - current
		if limit >= duration {
			res += duration
		} else {
			res += limit
		}

	}
	return res
}

// FindKthLargest 215. 数组中的第K个最大元素
func FindKthLargest(nums []int, k int) int {
	heapSize := len(nums)

	// 将数组转换成最大堆
	buildMaxHeap(nums, heapSize)

	for i := len(nums) - 1; i > len(nums)-k; i-- {
		// 移除当前堆顶到堆末尾
		nums[0], nums[i] = nums[i], nums[0]
		// 缩减堆尺寸, 数组中的顺序重新排序
		heapSize--
		maxHeapify(nums, 0, heapSize)
	}
	return nums[0]
}

// 构建一个堆, 遍历每一层结构
func buildMaxHeap(a []int, heapSize int) {
	for i := heapSize / 2; i >= 0; i-- {
		maxHeapify(a, i, heapSize)
	}
}

func maxHeapify(a []int, i, heapSize int) {
	largest, left, right := i, i*2+1, i*2+2

	// 左分支大于当前堆顶, 则标记左边为最大值
	if left < heapSize && a[left] > a[largest] {
		largest = left
	}

	// 右分支大于最大值, 则标记右边为最大值
	if right < heapSize && a[right] > a[largest] {
		largest = right
	}

	// 当前堆顶不止最大值, 则把最大值放到堆顶
	if largest != i {
		a[i], a[largest] = a[largest], a[i]
		maxHeapify(a, largest, heapSize)
	}
}
