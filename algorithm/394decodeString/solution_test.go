package _394decodeString

import (
	"github.com/astaxie/beego/logs"
	"testing"
)

func TestCase(t *testing.T) {
	demo := "3[a]2[bc]"
	expect := "aaabcbc"
	logs.Info("demo: %v expect: %v, res: %v", demo, expect, DecodeString(demo))
}

func TestCase2(t *testing.T) {
	demo := "3[a2[c]]"
	expect := "accaccacc"
	logs.Info("demo: %v expect: %v, res: %v", demo, expect, DecodeString(demo))
}

func TestCase3(t *testing.T) {
	demo := "2[abc]3[cd]ef"
	expect := "abcabccdcdcdef"
	logs.Info("demo: %v expect: %v, res: %v", demo, expect, DecodeString(demo))
}

func TestCase4(t *testing.T) {
	demo := "abc3[cd]xyz"
	expect := "abccdcdcdxyz"
	logs.Info("demo: %v expect: %v, res: %v", demo, expect, DecodeString(demo))
}
