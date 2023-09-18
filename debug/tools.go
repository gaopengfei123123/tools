package debug

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/gaopengfei123123/tools/convert"
	"os"
)

func init() {

}

// DD dump and die
func DD(f interface{}, v ...interface{}) {
	logs.Debug(f, v...)
	//b, _ := json.Marshal(convert.JsonCamelCase{v})
	//logs.Debug("%s %s", f, b)
	os.Exit(0)
}

func PrintJson(tag string, target interface{}, format ...bool) {
	// 如果日志模式是非打印模式, 则不再解析 json
	if logs.LevelDebug > logs.GetBeeLogger().GetLevel() {
		return
	}
	b, _ := convert.JSONEncode(target, format...)
	logs.Debug("%s: %s", tag, b)
}

// OutputJSON 输出json格式
func OutputJSON(k string, v ...interface{}) {
	b, _ := json.Marshal(v)
	logs.Debug("%s: %s", k, b)
}

func Info(f interface{}, v ...interface{}) {
	logs.Trace(f, v...)
}

func Debug(f interface{}, v ...interface{}) {
	logs.Debug(f, v...)
}
