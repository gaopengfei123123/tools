package _160getIntersectionNode

/*
*
  - Definition for singly-linked list.
  - type ListNode struct {
  - Val int
  - Next *ListNode
  - }

解题思路: 双指针

两个链表相交后数据内容相同, 因此如果两个链表相交, 那么他们从相交点算一定是对称的

指针遍历顺序如下
pa -> nodeA -> nodeB
pb -> nodeB -> nodeA

按 case1举例来说
nodeA := []int{4, 1, 8, 4, 5}
nodeB := []int{5, 6, 1, 8, 4, 5}

双指针拼装后
41845561845
56184541845
可以看到从8这个位置往后开始相同

详细解答链接 https://leetcode.cn/problems/intersection-of-two-linked-lists/solutions/10774/tu-jie-xiang-jiao-lian-biao-by-user7208t/?envType=study-plan-v2&envId=top-100-liked
*/
func getIntersectionNode(headA, headB *ListNode) *ListNode {
	if headA == nil || headB == nil {
		return nil
	}

	pa, pb := headA, headB
	for pa != pb {
		if pa == nil {
			pa = headB
		} else {
			pa = pa.Next
		}

		if pb == nil {
			pb = headA
		} else {
			pb = pb.Next
		}
	}

	return pa
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
