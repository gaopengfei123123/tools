package _84largestRectangleArea

// 解法链接: https://leetcode.cn/problems/largest-rectangle-in-histogram/description/?envType=study-plan-v2&envId=top-100-liked
/**
1. 使用单调栈
2. 所有入栈的高度一定是比上一个的高, 初始宽度都为1
3. 当出现当前高度比栈顶高度低的时候, 就以当前高度往栈内反着推, 借用栈内的宽度计算面积, 然后栈顶是当前高度, 宽度是当前次累计起来的
4. 注意出现边界情况,添加一个虚拟位置到最后, 高度为0
*/
func largestRectangleArea(heights []int) int {
	n := len(heights)

	stack := make([][2]int, 0, len(heights))

	ans := 0
	for i := 0; i <= n; i++ {
		height := 0
		if i < n {
			height = heights[i]
		}

		if len(stack) == 0 {
			stack = append(stack, [2]int{1, height})
			continue
		}

		// 如果出现更高的高度, 那么就插入最新的
		if height > stack[len(stack)-1][1] {
			stack = append(stack, [2]int{1, height})
			continue
		}

		// 如果出现相同高度的, 那么就延长宽度
		if height == stack[len(stack)-1][1] {
			stack[len(stack)-1][0] += 1
			continue
		}

		width := 0
		// **关键逻辑** 当出现一个更低的高度后, 开始计算当前栈内累计的最大面积
		for len(stack) > 0 && height < stack[len(stack)-1][1] {
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			width += top[0]
			ans = max(ans, width*top[1])
		}
		stack = append(stack, [2]int{width + 1, height})
	}
	return ans
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
