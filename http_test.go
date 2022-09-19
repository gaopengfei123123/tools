package tools

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"testing"
)

func TestPost(t *testing.T) {
	uri := "external/crop_wx/wt_room/room_info_by_name"
	host := "http://wukuapi.laidan.com"
	requestUrl := fmt.Sprintf("%s/%s", host, uri)
	params := map[string]string{
		"name":   "0401群1今日头条",
		"output": "json",
	}
	header := map[string]string{
		"name": "GPF",
	}
	body, err := HTTPPost(requestUrl, params, header)
	logs.Debug("response body: %s, err: %v", body, err)
}
