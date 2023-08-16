package convert

import (
	"github.com/astaxie/beego/logs"
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
