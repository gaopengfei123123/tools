package tools

import (
	"fmt"
	"strconv"
	"testing"
)

func TestDivide(t *testing.T) {
	t.Log("ttt")
	a := float32(2)
	b := 100
	num, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(a)/float64(b)), 64)
	t.Logf("res: %#+v", float32(num))
}

func TestArrayChunk(t *testing.T) {
	arr := []string{
		"a", "b", "3", "4", "5",
	}

	t.Logf("%v", ArrayChunkStr(arr, 2))
	t.Logf("%v", ArrayChunkStr(arr, 6))
}

func TestArrayChunk2(t *testing.T) {
	arr := []int32{
		1, 2, 3, 4, 5,
	}
	t.Logf("%v", ArrayChunkInt32(arr, 2))
	t.Logf("%v", ArrayChunkInt32(arr, 6))
}
