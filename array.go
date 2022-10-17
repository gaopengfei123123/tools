package tools

import (
	"math"
)

// ArrayChunkStr 数组分组
func ArrayChunkStr(array []string, num int) [][]string {
	max := len(array)
	if max <= num {
		return [][]string{array}
	}

	size := math.Ceil(float64(max) / float64(num))
	chunkNum := int(size)
	result := make([][]string, 0, chunkNum)
	for i := 0; i < max; i = i + num {
		end := i + num
		if end > max {
			end = max
		}
		result = append(result, array[i:end])
	}
	return result
}

// ArrayChunkInt32 数组分组
func ArrayChunkInt32(array []int32, num int) [][]int32 {
	max := len(array)
	if max <= num {
		return [][]int32{array}
	}

	size := math.Ceil(float64(max) / float64(num))
	chunkNum := int(size)
	result := make([][]int32, 0, chunkNum)
	for i := 0; i < max; i = i + num {
		end := i + num
		if end > max {
			end = max
		}
		result = append(result, array[i:end])
	}
	return result
}
