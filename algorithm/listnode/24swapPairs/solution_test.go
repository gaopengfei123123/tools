package _24swapPairs

import (
	"github.com/astaxie/beego/logs"
	"testing"
)

func TestCase(t *testing.T) {
	nodeA := []int{1, 2, 3, 4}
	expect := []int{2, 1, 4, 3}

	NodeA := BuildNodeList(nodeA)

	res := swapPairs(NodeA)
	logs.Info("nodeA: %v, expect: %v, res: %v", nodeA, expect, PrintNodeList(res))
}

func TestCase2(t *testing.T) {
	nodeA := []int{1}
	expect := []int{1}

	NodeA := BuildNodeList(nodeA)

	res := swapPairs(NodeA)
	logs.Info("nodeA: %v, expect: %v, res: %v", nodeA, expect, PrintNodeList(res))
}
