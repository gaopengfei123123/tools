package tools

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"testing"
)

func TestHTTPPost(t *testing.T) {
	uri := "/"
	host := "https://www.baidu.com"
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

func TestHTTPGet(t *testing.T) {
	uri := "/"
	host := "https://www.baidu.com"
	requestUrl := fmt.Sprintf("%s/%s", host, uri)
	params := map[string]string{
		"name":   "0401群1今日头条",
		"output": "json",
	}
	header := map[string]string{
		"name": "GPF",
	}
	body, err := HTTPGet(requestUrl, params, header)
	logs.Debug("response body: %s, err: %v", body, err)
}
