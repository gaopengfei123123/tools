package _141hasCycle

import (
	"github.com/astaxie/beego/logs"
	"testing"
)

func TestCase(t *testing.T) {
	nodeA := []int{1, 2, 3, 2, 1}

	NodeA := BuildNodeList(nodeA)

	expect := true
	logs.Info("nodeA: %v, expect: %v, res: %v", nodeA, expect, hasCycle(NodeA))
}

func TestCase2(t *testing.T) {
	nodeA := []int{1, 2}

	NodeA := BuildNodeList(nodeA)

	expect := false
	logs.Info("nodeA: %v, expect: %v, res: %v", nodeA, expect, hasCycle(NodeA))
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
