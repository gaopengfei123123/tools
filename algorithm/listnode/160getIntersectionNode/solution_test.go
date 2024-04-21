package _160getIntersectionNode

import (
	"github.com/astaxie/beego/logs"
	"testing"
)

func TestCase(t *testing.T) {
	nodeA := []int{4, 1, 8, 4, 5}
	nodeB := []int{5, 6, 1, 8, 4, 5}
	// 41845561845
	// 56184541845

	NodeA := BuildNodeList(nodeA)
	NodeB := BuildNodeList(nodeB)
	expect := 8

	res := getIntersectionNode(NodeA, NodeB)
	logs.Info("nodeA: %v, nodeB: %v,  expect: %v, res: %v", nodeA, nodeB, expect, res)
}

func TestCase2(t *testing.T) {
	//demo := []int{4, 2, 0, 3, 2, 5}
	//expect := 9
	//logs.Info("demo: %v expect: %v, res: %v", demo, expect, trap(demo))
}

func TestBuildNodeList(t *testing.T) {
	nodeA := []int{4, 1, 8, 4, 5}
	res := BuildNodeList(nodeA)
	logs.Info("build node: %v", res)
	logs.Info("build node: %v", res.Next)
	res = res.Next
	logs.Info("build node: %v", res.Next)
	res = res.Next
	logs.Info("build node: %v", res.Next)
	res = res.Next
	logs.Info("build node: %v", res.Next)
	res = res.Next
	logs.Info("build node: %v", res.Next)
	//logs.Info("Print node: %v", PrintNodeList(res))
}
