package _24swapPairs

/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */

/**

当前链表分解步骤如下如下:
head(current) -> A -> B -> tail
node1 = current.Next       =>   A -> B -> tail
node2 = current.Next.Next  =>   B -> tail
current.Next = node2       =>   head(current)  -> B -> tail
node1.Next = node2.Next    =>   A-> tail
node2.Next = node1         =>   B -> A -> tail
current = node1            =>   head -> B -> A (current) -> tail
*/

func swapPairs(head *ListNode) *ListNode {
	dummy := &ListNode{0, head}
	current := dummy
	for current.Next != nil && current.Next.Next != nil {
		node1 := current.Next
		node2 := current.Next.Next

		current.Next = node2

		node1.Next = node2.Next
		node2.Next = node1
		current = node1
	}
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
