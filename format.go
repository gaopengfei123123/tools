package tools

import (
	"github.com/astaxie/beego/logs"
	"github.com/gaopengfei123123/tools/convert"
	"time"
)

func TimeFormatInt32(tt int32, format ...string) (string, error) {
	tpl := "2006-01-02"
	if len(format) != 0 {
		tpl = format[0]
	}
	timer := time.Unix(int64(tt), 0).Format(tpl)
	return timer, nil
}

// PrintJson 将数据打印成 json格式
func PrintJson(tag string, data interface{}, format ...bool) {
	b, _ := convert.JSONEncode(data, format...)

	logs.Trace("%s: %s", tag, b)
}
