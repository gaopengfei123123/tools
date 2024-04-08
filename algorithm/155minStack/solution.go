package _155minStack

type MinStack struct {
	Stack    []int // 存数据的栈
	MinIndex []int // 存最小值下标的栈
}

func Constructor() MinStack {
	return MinStack{[]int{}, []int{}}
}

func (this *MinStack) Push(val int) {
	minIndex := len(this.Stack)
	this.Stack = append(this.Stack, val)

	if len(this.MinIndex) > 0 {
		curMin := this.Stack[this.MinIndex[len(this.MinIndex)-1]]
		if val <= curMin {
			this.MinIndex = append(this.MinIndex, minIndex)
		}
	} else {
		this.MinIndex = append(this.MinIndex, minIndex)
	}

}

func (this *MinStack) Pop() {
	curIndex := len(this.Stack) - 1
	this.Stack = this.Stack[:curIndex]

	if curIndex == this.MinIndex[len(this.MinIndex)-1] {
		this.MinIndex = this.MinIndex[:len(this.MinIndex)-1]
	}
}

func (this *MinStack) Top() int {
	if len(this.Stack) == 0 {
		return 0
	}
	return this.Stack[len(this.Stack)-1]
}

func (this *MinStack) GetMin() int {
	if len(this.MinIndex) == 0 {
		return 0
	}

	index := this.MinIndex[len(this.MinIndex)-1]

	return this.Stack[index]
}
