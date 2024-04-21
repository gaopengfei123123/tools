package topKFrequent

// TopKFrequentA 解法2
func TopKFrequentA(nums []int, k int) []int {
	// 1. 首先统计出各数字出现频率
	// 2. 将出现频率插入对应标号的桶中
	// 3. 返回对应的数字桶
	// 这种方法优点就是特别快, 缺点是得创建一个超大的空数组, 用空间换时间得典型
	var frequency map[int]int = make(map[int]int)
	for _, v := range nums {
		frequency[v] += 1
	}
	var bucket [][]int = make([][]int, len(nums)+1)
	var res []int
	for k, v := range frequency {
		bucket[v] = append(bucket[v], k)
	}
	for i, cnt := len(bucket)-1, 0; i >= 0 && cnt < k; i-- {
		res = append(res, bucket[i]...)
	}
	return res[:k]
}
