package _146LRUCache

import (
	"github.com/astaxie/beego/logs"
	"testing"
)

func TestCase(t *testing.T) {
	lru := Constructor(2)

	lru.Put(1, 1)
	lru.Put(2, 2)
	logs.Info("lru: %v", lru)

	logs.Info("get 1: %v, exist: %+v", lru.Get(1), lru)

	lru.Put(3, 3)
	logs.Info("Lru: %#+v", lru)
	////
	logs.Info("get -1 : %v,exist: %+v", lru.Get(2), lru)
	////
	lru.Put(4, 4)
	logs.Info("Lru: %#+v", lru)
	logs.Info("get 3: %v, exist: %+v", lru.Get(3), lru)
	logs.Info("get 4: %v, exist: %+v", lru.Get(4), lru)
}
