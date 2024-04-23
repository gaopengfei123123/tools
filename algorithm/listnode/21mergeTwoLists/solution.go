package _21mergeTwoLists

/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */
func mergeTwoLists(list1 *ListNode, list2 *ListNode) *ListNode {
	if list1 == nil && list2 == nil {
		return nil
	}

	newList := new(ListNode)
	cursor := newList

	for list1 != nil && list2 != nil {
		var tmpNode *ListNode

		if list1.Val > list2.Val {
			tmpNode = list2
			list2 = list2.Next
		} else {
			tmpNode = list1
			list1 = list1.Next
		}

		cursor.Next = tmpNode
		cursor = cursor.Next
	}

	if list1 != nil {
		cursor.Next = list1
	}
	if list2 != nil {
		cursor.Next = list2
	}

	return newList.Next
}

func minNode(node1, node2 *ListNode) *ListNode {
	if node2 == nil {
		return node1
	}

	if node1 == nil {
		return node2
	}
	if node1.Val > node2.Val {

		return node2
	}

	return node1
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
