package topKFrequent

// TopKFrequent  347 前 K 个高频元素
func TopKFrequent(nums []int, k int) []int {
	// 1. 首先把数字出现频次做一个哈希表
	// 2. 将频次按最大桶排序, 取前 k 个频率
	// 3. 在哈希表中找到对应的数字, 返回结果
	freqMap := make(map[int]int)
	freqArr := make([]int, 0, len(nums))
	for _, cur := range nums {
		if _, ok := freqMap[cur]; !ok {
			freqMap[cur] = 0
		}
		freqMap[cur]++
	}

	for _, v := range freqMap {
		freqArr = append(freqArr, v)
	}

	heapSize := len(freqArr)
	buildMaxHeaT(freqArr, heapSize)

	res := make([]int, 0, len(nums))
	for i := len(freqArr) - 1; i >= len(freqArr)-k; i-- {
		// 取当前排名最高的频率
		cur := freqArr[0]
		for v, freq := range freqMap {
			if freq == cur {
				res = append(res, v)
				delete(freqMap, v)
			}
		}
		freqArr[i], freqArr[0] = freqArr[0], freqArr[i]
		heapSize--
		maxHeapT(freqArr, 0, heapSize)
	}

	return res
}

func buildMaxHeaT(nums []int, heapSize int) {
	for i := heapSize / 2; i >= 0; i-- {
		maxHeapT(nums, i, heapSize)
	}
}

func maxHeapT(nums []int, i, heapSize int) {
	largest, left, right := i, i*2+1, i*2+2
	if left < heapSize && nums[left] > nums[largest] {
		largest = left
	}

	if right < heapSize && nums[right] > nums[largest] {
		largest = right
	}

	if largest != i {
		nums[i], nums[largest] = nums[largest], nums[i]
		maxHeapT(nums, largest, heapSize)
	}
}
