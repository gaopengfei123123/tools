package _141hasCycle

/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */

/**
解题思路:
1. 使用快慢指针对链表进行检测, 慢指针走一步, 快指针走两步, 如果存在环, 那么两个指针一定会相遇
2. 如果需要计算环的大小, 那么从第一次相遇开始计算, 到下一次相遇时遍历的次数就是环的大小
题解链接: https://leetcode.cn/problems/linked-list-cycle/?envType=study-plan-v2&envId=top-100-liked
*/

func hasCycle(head *ListNode) bool {
	if head == nil || head.Next == nil {
		return false
	}
	fast := head.Next
	slow := head

	for fast.Next != nil && fast.Next.Next != nil {
		if fast == slow {
			return true
		}
		fast = fast.Next.Next
		slow = slow.Next
	}
	return false
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
