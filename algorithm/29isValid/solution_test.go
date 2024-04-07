package _29isValid

import (
	"github.com/astaxie/beego/logs"
	"testing"
)

func TestIsValid(t *testing.T) {
	demo := "()"
	expect := true
	logs.Info("demo: %v expect: %v, res: %v", demo, expect, IsValid(demo))
}

func TestIsValid2(t *testing.T) {
	demo := "()[]{}"
	expect := true
	logs.Info("demo: %v expect: %v, res: %v", demo, expect, IsValid(demo))
}

func TestIsValid3(t *testing.T) {
	demo := "(]"
	expect := false
	logs.Info("demo: %v expect: %v, res: %v", demo, expect, IsValid(demo))
}

func TestIsValid4(t *testing.T) {
	demo := "]"
	expect := false
	logs.Info("demo: %v expect: %v, res: %v", demo, expect, IsValid(demo))
}
