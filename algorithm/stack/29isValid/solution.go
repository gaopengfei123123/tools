package _29isValid

/*
*
IsValid
利用栈的先进后出的特性, 栈最上面的括号一定是一层层闭合的
*/
func IsValid(s string) bool {
	// 允许的字符
	var allowMap = map[rune]struct{}{
		'(': {},
		')': {},
		'{': {},
		'}': {},
		'[': {},
		']': {},
	}

	// 可以抵消的字符
	var configMap = map[rune]rune{
		')': '(',
		'}': '{',
		']': '[',
	}

	stack := NewStack()
	arr := []rune(s)
	for _, v := range arr {
		//logs.Info("stack: %v, new:%v", string(stack), string(v))
		if _, exist := allowMap[v]; !exist {
			return false
		}

		if right, ok := configMap[v]; ok {
			if right == stack.Head() {
				stack.Pop()
				continue
			}
		}
		stack.Push(v)
	}
	//logs.Info("final stack: %v", string(stack))

	return stack.Len() == 0
}

type Stack []rune

func NewStack() Stack {
	return make([]rune, 0, 10)
}
func (st *Stack) Push(s rune) {
	*st = append(*st, s)
}

func (st *Stack) Pop() rune {
	if len(*st) == 0 {
		return 0
	}
	ln := len(*st)
	x := (*st)[ln-1]
	*st = (*st)[:ln-1]
	return x
}

func (st *Stack) Head() rune {
	if len(*st) == 0 {
		return 0
	}
	ln := len(*st)
	return (*st)[ln-1]
}

func (st *Stack) Len() int {
	return len(*st)
}
