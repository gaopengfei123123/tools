package _2addTwoNumbers

import (
	"github.com/astaxie/beego/logs"
	"testing"
)

func TestCase(t *testing.T) {
	nodeA := []int{2, 4, 3}
	nodeB := []int{5, 6, 4}

	NodeA := BuildNodeList(nodeA)
	NodeB := BuildNodeList(nodeB)

	expect := []int{7, 0, 8}
	res := addTwoNumbers(NodeA, NodeB)
	logs.Info("nodeA: %v, expect: %v, res: %v", nodeA, expect, PrintNodeList(res))
}

func TestCase2(t *testing.T) {
	nodeA := []int{0}
	nodeB := []int{0}

	NodeA := BuildNodeList(nodeA)
	NodeB := BuildNodeList(nodeB)

	expect := []int{0}
	res := addTwoNumbers(NodeA, NodeB)
	logs.Info("nodeA: %v, expect: %v, res: %v", nodeA, expect, PrintNodeList(res))
}

func TestCase3(t *testing.T) {
	nodeA := []int{9, 9, 9, 9, 9, 9, 9}
	nodeB := []int{9, 9, 9, 9}

	NodeA := BuildNodeList(nodeA)
	NodeB := BuildNodeList(nodeB)

	expect := []int{8, 9, 9, 9, 0, 0, 0, 1}
	res := addTwoNumbers(NodeA, NodeB)
	logs.Info("nodeA: %v, expect: %v, res: %v", nodeA, expect, PrintNodeList(res))
}
