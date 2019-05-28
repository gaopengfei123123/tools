package logs

import (
	"encoding/json"
	"os"
	"runtime"
	"time"
)

func init() {
	Register(AdapterConsole, NewConsole)
}

type consoleWriter struct {
	lg       *logWriter
	Level    int  `json:"level"`
	Colorful bool `json:"color"` // 如果terminal支持的话
}

// NewConsole 返回日志实体
func NewConsole() LoggerInterface {
	return &consoleWriter{
		lg:       newLogWriter(os.Stdout),
		Level:    LevelDebug,
		Colorful: runtime.GOOS != "windows",
	}
}

func (c *consoleWriter) Init(jsonConfig string) error {
	if len(jsonConfig) == 0 {
		return nil
	}
	err := json.Unmarshal([]byte(jsonConfig), c)
	if runtime.GOOS == "windows" {
		c.Colorful = false
	}
	return err
}

func (c *consoleWriter) WriteMsg(when time.Time, msg string, level int) error {
	if level > c.Level {
		return nil
	}

	if c.Colorful {
		msg = colors[level](msg)
	}
	c.lg.println(when, msg)
	return nil
}

func (c *consoleWriter) Destroy() {

}
func (c *consoleWriter) Flush() {

}

// 对颜色的管理
type brush func(string) string

func newBrush(color string) brush {
	pre := "\033["
	reset := "\033[0m"
	return func(text string) string {
		return pre + color + "m" + text + reset
	}
}

// 颜色设置以及字体设置参考 https://blog.csdn.net/guohaosun/article/details/83032791
var colors = []brush{
	newBrush("1;37"), // Emergency          white
	newBrush("1;36"), // Alert              cyan
	newBrush("1;35"), // Critical           magenta
	newBrush("1;31"), // Error              red
	newBrush("1;33"), // Warning            yellow
	newBrush("1;32"), // Notice             green
	newBrush("1;34"), // Informational      blue
	newBrush("7;33"), // Debug              Background yellow
}
