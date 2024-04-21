package _206reverseList

/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */
func reverseList(head *ListNode) *ListNode {
	var pre *ListNode
	cur := head

	for cur != nil {
		tmp := cur.Next // 缓存当前节点后续的数据
		cur.Next = pre  // 将当前节点的后续指针指向新的链表

		pre = cur // 将当前节点迁移到新的列表中, 此时新的链表是反着加的, 始终指向前一个节点
		cur = tmp // 将数据窗口划向老链表的下一个节点
	}

	return pre
}

type ListNode struct {
	Val  int
	Next *ListNode
}

func BuildNodeList(demo []int) *ListNode {
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

	for {
		res = append(res, node.Val)

		if node.Next == nil {
			break
		}
		node = node.Next
	}
	return res
}
