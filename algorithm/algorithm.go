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

func FindKthLargest(nums []int, k int) int {
	//nums = _quickSort(nums)

	return 0
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
