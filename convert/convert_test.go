package convert

import "testing"

func TestDiv(t *testing.T) {
	for i := 1; i < 30; i++ {
		t.Logf("row:%v => %v\n", i, Div(i))
	}
}
