package _42trap

/*
*

接雨水问题, 解题方法和 84求最大面积思路一样
1. 使用递减栈, 栈内至少需要存每个高度的下标, 栈内每个元素的高度是依次递减的
2. 当出现一个高度高于栈顶元素的高度时, 说明碗的右侧找到了, 这时候找碗的左侧高度, 取两者最小值, 算是碗的容积高度
3. 每当计算完一层栈后要抛出, 不重复计算面积,
示例数组:[2,1,0,1,3]

当偏移下标是3时, 栈内是 [2,1,0], 右边是1, 当前碗最高高度是min(1,1)=1, 累计承接雨滴 1
当偏移下标是4是, 栈内是 [2,1,1,1] 右边是3, 碗最高高度是 [min(2,3)-1]=1, 可承接雨滴是3

变化值为:
[2,1,[0],1,3]  计算到初始的那个小坑, 容积是1
[2,[1,1,1],3]  计算到较大的那个坑, 容积是3

示例图参考: https://leetcode.cn/problems/trapping-rain-water/solutions/616404/42-jie-yu-shui-shuang-zhi-zhen-dong-tai-wguic/
*/
func trap(heights []int) int {
	stack := make([][2]int, 0)

	water := 0
	n := len(heights)
	for i := 0; i < n; i++ {
		height := heights[i]

		if len(stack) == 0 {
			stack = append(stack, [2]int{i, height})
			continue
		}

		// 如果出现当前高度大于栈顶高度, 则反向计算可乘的水滴数
		for len(stack) > 0 && stack[len(stack)-1][1] < height {
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if len(stack) == 0 {
				break
			}

			left := stack[len(stack)-1]  // 获取碗的左边高度
			edge := min(left[1], height) // 确定碗最终高度是多少
			h := edge - top[1]           // 计算当前格的可承接雨水的量
			w := i - left[0] - 1         // 计算碗的宽度
			water += w * h
		}

		stack = append(stack, [2]int{i, height})
	}
	return water
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
