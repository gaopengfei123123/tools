package _19addTwoNumbers

import (
	"github.com/astaxie/beego/logs"
)

/*
*
解题思路:
1. 设置一个快慢指针, 并且设置一个哨兵节点
2. 快指针和慢指针间隔 n 个节点, 当快指针遍历完成后, 慢指针就处在要移除的节点前一个位置
3. 设置哨兵节点的目的是防止出现空值这类的边界判断
*/
func removeNthFromEnd(head *ListNode, n int) *ListNode {
	fast := head
	dummy := &ListNode{0, head}
	slow := dummy

	step := 0
	for fast.Next != nil {
		fast = fast.Next

		step++
		if step >= n {
			slow = slow.Next
		}

	}
	logs.Info("slow: %v, fast: %v", slow, fast)

	// 把慢指针的下一节移除掉
	slow.Next = slow.Next.Next

	return dummy.Next
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
