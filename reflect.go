package tools

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"reflect"
)

// StructByReflect 遍历struct并且自动进行赋值
func StructByReflect(beforeMap map[string]interface{}, inStructPtr interface{}) error {
	marshal, err := json.Marshal(beforeMap)
	if err != nil {
		return err
	}

	err = json.Unmarshal(marshal, inStructPtr)
	if err != nil {
		return err
	}
	return nil
}

type KT struct {
	Key  string
	Type string
}

// GetStructKeyType 获取struct的kv结构
func GetStructKeyType(structName interface{}) ([]KT, error) {
	t := reflect.TypeOf(structName)
	//rv := reflect.ValueOf(structName)
	if t.Kind() != reflect.Struct {
		err := errors.Errorf("%s is not struct", t.Name())
		return nil, err
	}

	result := make([]KT, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		//debug.Info("key: %#+v type: %v  value: %#+v", t.Field(i).Name, t.Field(i).Type, rv.Field(i).Interface())
		result[i] = KT{
			Key:  t.Field(i).Name,
			Type: fmt.Sprintf("%v", t.Field(i).Type),
		}
	}
	return result, nil
}

// GetStructStringField 获取struct中指定key的string值
func GetStructStringField(input interface{}, key string) (value string, err error) {
	v, err := getStructField(input, key)
	if err != nil {
		return
	}

	value, ok := v.(string)
	if !ok {
		return value, errors.New("can't convert key'v to string")
	}

	return
}

func getStructField(input interface{}, key string) (value interface{}, err error) {
	rv := reflect.ValueOf(input)
	rt := reflect.TypeOf(input)
	if rt.Kind() != reflect.Struct {
		return value, errors.New("input must be struct")
	}

	keyExist := false
	for i := 0; i < rt.NumField(); i++ {
		curField := rv.Field(i)
		if rt.Field(i).Name == key {
			switch curField.Kind() {
			case reflect.String, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int, reflect.Float64, reflect.Float32:
				keyExist = true
				value = curField.Interface()
			default:
				return value, errors.New("key must be int float or string")
			}
		}
	}

	if !keyExist {
		return value, errors.New(fmt.Sprintf("key %s not found in %s's field", key, rt))
	}

	return
}

// 分解单个参数, 将值动态的赋给索引, 如果类型不一致, 则会报错
func resultUnmarshal(src interface{}, dst interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("字段赋值错误: 源字段类型:%v, 赋值字段是%v", reflect.ValueOf(dst).Elem().Type(), reflect.ValueOf(src).Elem().Type())
			return
		}
	}()

	dstRv := reflect.ValueOf(dst)
	if dstRv.Kind() != reflect.Ptr || dstRv.IsNil() {
		return fmt.Errorf("params[%v]不是一个引用类型参数", dst)
	}
	// 源是空的时候不做处理
	if src == nil {
		return nil
	}

	// 利用反射给字段赋值
	srcRv := reflect.ValueOf(src)
	//logs.Trace("rv dst: %v, src: %v", dstRv.Elem().Type(), srcRv.Type())
	dstRv.Elem().Set(srcRv)
	return nil
}

// InterfaceToResult 将interface 里面的字段赋值给后面各种变量
func InterfaceToResult(resultList []interface{}, returnItems ...interface{}) error {
	allowIndex := len(resultList)
	//logs.Trace("InterfaceToResult resultList: %#+v, len: %v", resultList, allowIndex)

	for i := 0; i < len(returnItems); i++ {
		if i >= allowIndex {
			continue
		}
		//logs.Trace("i: %v src: %v, dst: %v", i, resultList[i], returnItems[i])
		err := resultUnmarshal(resultList[i], returnItems[i])
		if err != nil {
			return err
		}
	}
	return nil
}
