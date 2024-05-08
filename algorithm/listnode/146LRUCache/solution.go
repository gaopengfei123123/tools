package _146LRUCache

/**
解题思路:
1. 为保证Put Get操作时间复杂度是 O(1), 那么需要维护一个双向链表, 掐住头和尾, 同时, 链表的头和尾就是固定的两个值不要动了, 也是为了方便维护
2. 把对链表的变更动作独立出来,方便理清思路, 简单动作 addhead, removeNode, 复合动作 removeToHead, removeTail

链接: https://leetcode.cn/problems/lru-cache/solutions/259678/lruhuan-cun-ji-zhi-by-leetcode-solution/?envType=study-plan-v2&envId=top-100-liked
*/

type LRUCache struct {
	Cap     int
	Size    int
	Head    *DLinkedNode
	Tail    *DLinkedNode
	HashMap map[int]*DLinkedNode
}

func Constructor(capacity int) LRUCache {
	lru := LRUCache{
		Cap:     capacity,
		Head:    InitNode(0, 0),
		Tail:    InitNode(0, 0),
		HashMap: make(map[int]*DLinkedNode),
	}

	lru.Head.Next = lru.Tail
	lru.Tail.Pre = lru.Head

	return lru
}

func InitNode(key, value int) *DLinkedNode {
	return &DLinkedNode{
		Val: value,
		Key: key,
	}
}

func (this *LRUCache) Get(key int) int {
	if _, ok := this.HashMap[key]; !ok {
		return -1
	}
	node := this.HashMap[key]
	this.removeToHead(node)
	return node.Val
}

func (this *LRUCache) Put(key int, value int) {
	if _, exist := this.HashMap[key]; !exist {
		node := InitNode(key, value)
		this.addHead(node)
		this.HashMap[key] = node
		this.Size++

		for this.Size > this.Cap {
			node = this.removeTail()
			delete(this.HashMap, node.Key)
			this.Size--
		}

	} else {
		node := this.HashMap[key]
		node.Val = value
		this.removeToHead(node)
	}
}

func (this *LRUCache) addHead(node *DLinkedNode) {
	node.Pre = this.Head
	node.Next = this.Head.Next
	this.Head.Next.Pre = node
	this.Head.Next = node
}

func (this *LRUCache) removeNode(node *DLinkedNode) {
	node.Pre.Next = node.Next
	node.Next.Pre = node.Pre
}

func (this *LRUCache) removeToHead(node *DLinkedNode) {
	this.removeNode(node)
	this.addHead(node)
}

func (this *LRUCache) removeTail() *DLinkedNode {
	node := this.Tail.Pre
	this.removeNode(node)
	return node
}

/**
 * Your LRUCache object will be instantiated and called as such:
 * obj := Constructor(capacity);
 * param_1 := obj.Get(key);
 * obj.Put(key,value);
 */

type DLinkedNode struct {
	Val  int
	Key  int
	Pre  *DLinkedNode
	Next *DLinkedNode
}

func BuildNodeList(demo []int) *DLinkedNode {
	if len(demo) == 0 {
		return nil
	}
	res := new(DLinkedNode)

	for i := len(demo) - 1; i >= 0; i-- {
		cur := demo[i]
		//logs.Info("cur:%v", cur)
		res.Val = cur

		tmp := new(DLinkedNode)
		tmp.Next = res
		res = tmp
	}
	return res.Next
}

func PrintNodeList(node *DLinkedNode) []int {
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
