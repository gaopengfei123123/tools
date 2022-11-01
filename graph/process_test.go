package graph

import (
	"testing"
	"time"
)

// 进度条实现 https://segmentfault.com/a/1190000023375330
func TestBar_NewOption(t *testing.T) {
	var bar Bar
	bar.NewOption(0, 100)
	for i := 0; i <= 100; i++ {
		time.Sleep(100 * time.Millisecond)
		bar.Play(int64(i))
	}
	bar.Finish()
}
