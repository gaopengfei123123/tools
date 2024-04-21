package _206reverseList

import (
	"github.com/astaxie/beego/logs"
	"testing"
)

func TestCase(t *testing.T) {
	nodeA := []int{1, 2, 3, 4, 5}
	nodeB := []int{5, 4, 3, 2, 1}

	NodeA := BuildNodeList(nodeA)
	NodeB := BuildNodeList(nodeB)

	res := reverseList(NodeA)
	logs.Info("nodeA: %v,  expect: %v, res: %v", nodeA, PrintNodeList(NodeB), PrintNodeList(res))
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
