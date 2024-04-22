package _142detectCycle

/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */

/*
*
解题思路:
1. 使用快慢指针对链表进行检测, 慢指针走一步, 快指针走两步, 如果存在环, 那么两个指针一定会相遇
2. 并且如果快指针速度是2, 慢指针速度是1, 那么快慢指针一定在快指针走了一圈以后相遇, 因为如果是一个纯环形, 那么当慢指针走一半时,快指针一定是套了快指针一圈
因此, 假设从开头到环起始点的距离是 a, 慢指针从环起点到相遇时走的路程是 b, 剩下的距离到环起点是 c, 有整个环长度是 b+c

那么又因为快指针走的长度是慢指针的2倍, 因此有关系式:

2(a+b) = a+b+(b+c)

简化得: a = c

因此, 当快指针和慢指针第一次相遇时, 这时候从链表的起点设置一个和慢指针相同速度的 tmp指针, 那么 tmp指针将和慢指针在环的起始点相遇

题解链接: https://leetcode.cn/problems/linked-list-cycle-ii/solutions/2751456/kuai-man-zhi-zhen-fang-shi-zhao-huan-xin-d4zt/
*/
func detectCycle(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return nil
	}
	slow := head
	fast := head

	for fast != nil && fast.Next != nil {
		fast = fast.Next.Next
		slow = slow.Next

		if slow == fast {
			tmp := head
			for tmp != slow {
				tmp = tmp.Next
				slow = slow.Next
			}
			return tmp
		}
	}

	return nil
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
