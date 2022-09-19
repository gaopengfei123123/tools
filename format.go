package tools

import "time"

func TimeFormatInt32(tt int32, format ...string) (string, error) {
	tpl := "2006-01-02"
	if len(format) != 0 {
		tpl = format[0]
	}
	timer := time.Unix(int64(tt), 0).Format(tpl)
	return timer, nil
}
