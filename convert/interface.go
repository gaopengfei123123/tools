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

	// 目前有个结构体嵌套的问题, 先不处理
	//// 获取赋值的字段名和偏移量字典, 规避下n^2的问题
	//bValMap := make(map[string]int)
	//for i := 0; i < bVal.NumField(); i++ {
	//	if bVal.Field(i).IsValid() && bVal.Field(i).CanSet() {
	//		name := bVal.Type().Field(i).Name
	//		bValMap[name] = i
	//	}
	//}

	for i := 0; i < vVal.NumField(); i++ {
		// 在要修改的结构体中查询有数据结构体中相同属性的字段，有则修改其值
		name := vTypeOfT.Field(i).Name

		//if index, ok := bValMap[name]; ok {
		//	bVal.Field(index).Set(reflect.ValueOf(vVal.Field(i).Interface()))
		//}

		if bVal.FieldByName(name).IsValid() && bVal.FieldByName(name).CanSet() {
			bVal.FieldByName(name).Set(reflect.ValueOf(vVal.Field(i).Interface()))
		}
	}
}
