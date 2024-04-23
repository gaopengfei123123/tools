package _2addTwoNumbers

/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */
func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	result := new(ListNode)

	cursor := result
	tmpVal := 0

	for l1 != nil || l2 != nil {
		unit, decade := addNode(l1, l2, tmpVal)
		tmpVal = decade

		tmp := new(ListNode)
		tmp.Val = unit

		cursor.Next = tmp
		cursor = cursor.Next

		if l1 != nil {
			l1 = l1.Next
		}

		if l2 != nil {
			l2 = l2.Next
		}
	}

	// 将最高位放到链表尾部
	if tmpVal != 0 {
		tmp := new(ListNode)
		tmp.Val = tmpVal
		cursor.Next = tmp
	}

	return result.Next
}

func addNode(node1, node2 *ListNode, extra int) (unit, decade int) {
	var v1, v2 int
	if node1 != nil {
		v1 = node1.Val
	}
	if node2 != nil {
		v2 = node2.Val
	}
	unit = v1 + v2 + extra

	if unit >= 10 {
		decade = unit / 10
		unit = unit % 10
	}
	return
}

type ListNode struct {
	Val  int
	Next *ListNode
}

func BuildNodeList(demo []int) *ListNode {
	if len(demo) == 0 {
		return nil
	}
	res := new(ListNode)

	for i := len(demo) - 1; i >= 0; i-- {
		cur := demo[i]
		//logs.Info("cur:%v", cur)
		res.Val = cur

		tmp := new(ListNode)
		tmp.Next = res
		res = tmp
	}
	return res.Next
}

func PrintNodeList(node *ListNode) []int {
	res := make([]int, 0, 4)

	if node == nil {
		return nil
	}

	for {
		res = append(res, node.Val)

		if node.Next == nil {
			break
		}
		node = node.Next
	}
	return res
}
