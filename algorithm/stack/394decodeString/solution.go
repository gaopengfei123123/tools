package _394decodeString

func DecodeString(str string) string {
	strStack := make([]string, 0, 1)
	countStack := make([]int, 0, 1)
	curStr := ""  // 当前处理中的字符串
	curCount := 0 // 当前数字

	for _, char := range str {
		switch char {
		// 开始
		case '[':
			// 将最近的倍数放到栈中
			countStack = append(countStack, curCount)
			curCount = 0
			// 将上一次处理后的字符串放入栈中
			strStack = append(strStack, curStr)
			curStr = ""
			// 结束
		case ']':
			// 获取最近的数字, 然后复制当前缓存字符串
			preCount := countStack[len(countStack)-1]
			countStack = countStack[:len(countStack)-1]

			tmpStr := ""
			// 重复字符倍数
			for i := 0; i < preCount; i++ {
				tmpStr += curStr
			}

			// 获取上一层已经处理完的字符串, 再把当前[]下处理好的字符串拼上去
			preStr := strStack[len(strStack)-1]
			strStack = strStack[:len(strStack)-1]
			curStr = preStr + tmpStr
		default:
			// 存储数字
			num := int(char - '0')
			if num >= 0 && num <= 9 {
				curCount = curCount*10 + num
			} else {
				// 非关键词的放入到待处理字符串中
				curStr += string(char)
			}
		}
	}
	return curStr
}
