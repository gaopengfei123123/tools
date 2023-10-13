package convert

import (
	"github.com/astaxie/beego/logs"
	"reflect"
	"strconv"
)

// GetInt32ValueInMap 提取数据
func GetInt32ValueInMap(mapData map[string]interface{}, key string) int32 {
	v, ok := mapData[key]
	if !ok {
		return 0
	}
	switch v.(type) {
	case nil:
		return 0
	case int64:
		return int32(v.(int64))
	case float64:
		return int32(v.(float64))
	case string:
		n, err := strconv.Atoi(v.(string))
		if err != nil {
			return 0
		}
		return int32(n)
	case int32:
		return v.(int32)
	default:
		logs.Debug("mapData type unknown\n")
	}
	return 0
}

func GetFloat32ValueInMap(mapData map[string]interface{}, key string) float32 {
	v, ok := mapData[key]
	if !ok {
		return 0
	}
	switch v.(type) {
	case nil:
		return 0
	case int64:
		return float32(v.(int64))
	case float64:
		return float32(v.(float64))
	case string:
		n, err := strconv.ParseFloat(v.(string), 32)
		if err != nil {
			return 0
		}
		return float32(n)
	case int32:
		return float32(v.(int32))
	default:
		logs.Debug("mapData type unknown\n")
	}
	return 0
}

func GetInt64ValueInMap(mapData map[string]interface{}, key string) int64 {
	v, ok := mapData[key]
	if !ok {
		return 0
	}
	switch v.(type) {
	case nil:
		return 0
	case int64:
		return v.(int64)
	case float64:
		return int64(v.(float64))
	case string:
		n, err := strconv.Atoi(v.(string))
		if err != nil {
			return 0
		}
		return int64(n)
	case int32:
		return int64(v.(int32))
	default:
		logs.Debug("mapData type unknown\n")
	}
	return int64(0)
}

// InterfaceToInt32  转int32
func InterfaceToInt32(v interface{}) int32 {
	switch v.(type) {
	case nil:
		return 0
	case int64:
		return int32(v.(int64))
	case float64:
		return int32(v.(float64))
	case string:
		n, err := strconv.Atoi(v.(string))
		if err != nil {
			return 0
		}
		return int32(n)
	case int32:
		return v.(int32)
	default:
		logs.Debug("mapData type unknown\n")
	}
	return 0
}

// StructAssign  将value 中的值赋值到 bingding 中, 需要字段名, 字段类型一致
// binding type interface 要修改的结构体
// value type interace 有数据的结构体
func StructAssign(binding interface{}, value interface{}) {
	bVal := reflect.ValueOf(binding).Elem() //获取reflect.Type类型
	vVal := reflect.ValueOf(value).Elem()   //获取reflect.Type类型
	vTypeOfT := vVal.Type()
	for i := 0; i < vVal.NumField(); i++ {
		// 在要修改的结构体中查询有数据结构体中相同属性的字段，有则修改其值
		name := vTypeOfT.Field(i).Name
		if ok := bVal.FieldByName(name).IsValid(); ok {
			bVal.FieldByName(name).Set(reflect.ValueOf(vVal.Field(i).Interface()))
		}
	}
}
