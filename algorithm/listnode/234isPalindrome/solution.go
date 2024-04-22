package _234isPalindrome

/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */

/**
解题思路:
1. 先使用快慢指针, 来确定链表的中间节点是哪里
2. 反转后半部分链表, 并和前半部分比较是否为回文
3. 比较完成后, 将后半部分链表再反转回来

这种方法的缺点就是会内部改变链表结构, 如果涉及到多线程的话, 需要对资源上锁

题解链接: https://leetcode.cn/problems/palindrome-linked-list/solutions/457059/hui-wen-lian-biao-by-leetcode-solution/?envType=study-plan-v2&envId=top-100-liked
*/

func isPalindrome(head *ListNode) bool {
	// 1. 快慢指针确定中间节点
	rightHalfList := findListMiddle(head)

	// 2. 反转右侧链表数据
	rightReverse := nodeListReverse(rightHalfList.Next)

	// 3. 双指针比对是否回文
	pLeft := head
	pRight := rightReverse

	result := true
	for result && pRight != nil {
		if pLeft.Val != pRight.Val {
			result = false
		}
		pLeft = pLeft.Next
		pRight = pRight.Next
	}

	// 4. 将右侧反转的链表恢复
	rightHalfList.Next = nodeListReverse(rightReverse)

	return result
}

func nodeListReverse(head *ListNode) *ListNode {
	var pre *ListNode
	cur := head

	for cur != nil {
		tmp := cur.Next
		cur.Next = pre

		pre = cur
		cur = tmp
	}
	return pre
}

func findListMiddle(head *ListNode) *ListNode {
	slow := head
	fast := head

	for fast.Next != nil && fast.Next.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
	}

	return slow
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
