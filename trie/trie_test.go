package trie

import (
	"github.com/astaxie/beego/logs"
	"github.com/gaopengfei123123/tools/debug"
	"testing"
)

func TestTrie_SignSensitiveWords(t *testing.T) {
	rule := new(Trie)
	rule.BlackList = []string{"逆旅", "行人", "相见", "有心人", "人"}
	rule.BlackListToTrieTree()

	rawText := "人生如逆旅,我亦是行人,但愿初相见,不负有心人"

	highLightStr, hitKey := rule.SignSensitiveWords(rawText, "<tag>", "</tag>")
	logs.Info("\n处理前: %s\n处理后: %s", rawText, highLightStr)
	debug.PrintJson("命中关键词", hitKey)
}
