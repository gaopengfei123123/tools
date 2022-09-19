package tools

import "testing"

func TestGetCurrPath(t *testing.T) {
	p := GetCurrPath()
	t.Log(p)
}
