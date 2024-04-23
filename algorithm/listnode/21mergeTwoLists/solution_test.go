package _21mergeTwoLists

import (
	"github.com/astaxie/beego/logs"
	"testing"
)

func TestCase(t *testing.T) {
	nodeA := []int{1, 3, 5}
	nodeB := []int{2, 4, 6}

	NodeA := BuildNodeList(nodeA)
	NodeB := BuildNodeList(nodeB)

	expect := []int{1, 2, 3, 4, 5, 6}
	res := mergeTwoLists(NodeA, NodeB)
	logs.Info("nodeA: %v, expect: %v, res: %v", nodeA, expect, PrintNodeList(res))
}

func TestCase2(t *testing.T) {
	nodeA := []int{}
	nodeB := []int{2, 4, 6}

	NodeA := BuildNodeList(nodeA)
	NodeB := BuildNodeList(nodeB)

	expect := []int{2, 4, 6}
	res := mergeTwoLists(NodeA, NodeB)
	logs.Info("nodeA: %v, expect: %v, res: %v", nodeA, expect, PrintNodeList(res))
}

func TestCase3(t *testing.T) {
	nodeA := []int{}
	nodeB := []int{}

	NodeA := BuildNodeList(nodeA)
	NodeB := BuildNodeList(nodeB)

	expect := []int{}
	res := mergeTwoLists(NodeA, NodeB)
	logs.Info("nodeA: %v, expect: %v, res: %v", nodeA, expect, PrintNodeList(res))
}
