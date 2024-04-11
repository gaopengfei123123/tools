package _739dailyTemperatures

func dailyTemperatures(temperatures []int) []int {
	dayStack := make([]int, 0) // 设置一个递减栈, 栈底值是最大的
	res := make([]int, len(temperatures))
	for i := 0; i < len(temperatures); i++ {
		curTemp := temperatures[i]
		// 如果出现一个值比栈顶值大, 那么回溯整个栈, 栈内所有下标都和当前值比较大小, 直到当前值小于栈顶值
		for len(dayStack) > 0 && curTemp > temperatures[dayStack[len(dayStack)-1]] {
			preIndex := dayStack[len(dayStack)-1]
			dayStack = dayStack[:len(dayStack)-1]
			res[preIndex] = i - preIndex
		}
		// 放入栈中的是数组下标, 不是数组值
		dayStack = append(dayStack, i)
	}

	return res
}
