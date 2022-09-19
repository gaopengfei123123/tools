package tools

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strings"
)

// HTTPGet 简化请求
func HTTPGet(requestURI string, params map[string]string, header map[string]string) ([]byte, error) {
	query := ""
	if params != nil {
		for k, v := range params {
			query += fmt.Sprintf("&%s=%s", k, v)
		}
	}
	//logs.Debug("requestURI: %s", requestURI)
	requestURL := fmt.Sprintf("%s?%s", requestURI, query)
	//logs.Debug("requestURI: %s", requestURL)
	req, _ := http.NewRequest("GET", requestURL, nil)

	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}

	resp, err := (&http.Client{}).Do(req)
	defer resp.Body.Close()
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("%s curl fail", requestURI))
	}
	return ioutil.ReadAll(resp.Body)
}

func HTTPPost(requestURI string, params map[string]string, header map[string]string) ([]byte, error) {
	strArr := []string{}
	for k, v := range params {
		strArr = append(strArr, fmt.Sprintf("%v=%v", k, v))
	}
	strList := strings.NewReader(strings.Join(strArr, "&"))

	req, _ := http.NewRequest("POST", requestURI, strList)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}

	resp, err := (&http.Client{}).Do(req)
	defer resp.Body.Close()
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("%s curl fail", requestURI))
	}
	return ioutil.ReadAll(resp.Body)
}
