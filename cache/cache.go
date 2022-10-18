package cache

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"github.com/gaopengfei123123/tools"
	"github.com/pkg/errors"
	"reflect"
	"runtime"
	"strings"
)

// Encode 进行 golang 的序列化
func Encode(v interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(v); err != nil {
		return buf.Bytes(), err
	}
	return buf.Bytes(), nil
}

// Decode 将字节解析成对象
func Decode(v []byte, target interface{}) error {
	ior := bytes.NewReader(v)
	d := gob.NewDecoder(ior)
	err := d.Decode(target)
	return err
}

// CallFuncBody 调用函数的结构体
type CallFuncBody struct {
	cache    CommonDrive
	FuncName interface{} // 这里存放函数实体
	Params   []interface{}
	Result   []interface{}
	Err      error
}

func (cb *CallFuncBody) GetResult(returnItems ...interface{}) error {
	if cb.Err != nil {
		return cb.Err
	}

	return tools.InterfaceToResult(cb.Result, returnItems...)
}

func (cb *CallFuncBody) GetCacheKey() (key string, funcName string, err error) {
	paramsStr := fmt.Sprintf("%v", cb.Params)
	h := md5.New()
	h.Write([]byte(paramsStr))
	paramsStr = hex.EncodeToString(h.Sum(nil))
	//logs.Info("%v", GetFuncName(cb.FuncName))
	funcName = GetFuncName(cb.FuncName)
	key = fmt.Sprintf("CacheFuncKey:%s:%v", GetFuncName(cb.FuncName), paramsStr)
	return
}

// GetFuncName 这里取函数最后的名字
func GetFuncName(fc interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(fc).Pointer()).Name()
	strArr := strings.Split(name, "/")
	index := len(strArr) - 1
	lastName := strArr[index]
	return lastName
}

// CallFunc 利用反射动态执行函数, 直接从 batchExec 那边搬过来的
func CallFunc(body CallFuncBody) (result []interface{}, err error) {
	// 校验是否是函数
	if reflect.TypeOf(body.FuncName).Kind() != reflect.Func {
		err = errors.New(fmt.Sprintf("this is not a  func name abort. FuncName: %v", body.FuncName))
		return
	}
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = fmt.Errorf("%v", panicErr)
		}
	}()

	// 执行方法
	f := reflect.ValueOf(body.FuncName)
	// 校验传参值数量
	if len(body.Params) != f.Type().NumIn() {
		err = errors.New("The number of params is not adapted.")
		return
	}
	in := make([]reflect.Value, len(body.Params))
	for k, param := range body.Params {
		in[k] = reflect.ValueOf(param)
	}
	res := f.Call(in)
	result = make([]interface{}, len(res))
	for k, v := range res {
		result[k] = v.Interface()
	}
	return
}
