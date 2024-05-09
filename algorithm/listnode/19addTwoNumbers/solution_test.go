package _19addTwoNumbers

import "testing"

import (
	"github.com/astaxie/beego/logs"
)

func TestCase(t *testing.T) {
	nodeA := []int{1, 2, 3, 4, 5}
	expect := []int{1, 2, 3, 5}

	NodeA := BuildNodeList(nodeA)

	n := 2
	res := removeNthFromEnd(NodeA, n)
	logs.Info("nodeA: %v, expect: %v, res: %v", nodeA, expect, PrintNodeList(res))
}

func TestCase2(t *testing.T) {
	nodeA := []int{1}
	expect := []int{}

	NodeA := BuildNodeList(nodeA)

	n := 1
	res := removeNthFromEnd(NodeA, n)
	logs.Info("nodeA: %v, expect: %v, res: %v", nodeA, expect, PrintNodeList(res))
}
