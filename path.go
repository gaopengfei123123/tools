package tools

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GetCurrPath 获取当前可执行文件所在的磁盘路径
func GetCurrPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))
	ret := path[:index]
	return ret
}
