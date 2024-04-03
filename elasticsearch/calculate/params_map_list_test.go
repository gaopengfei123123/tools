package calculate

import (
	"github.com/gaopengfei123123/tools/debug"
	"testing"
)

// 将参数转换成query 基础示例
func TestParamsMapList_GenerateQuery(t *testing.T) {
	params := map[string]interface{}{
		"operate_config_id":                     1654,
		"small_course_order.small_course_id":    7094,
		"small_course_order.small_course_stage": 1,
	}

	boolQuery := new(ParamsMapList).LoadConfig("", params).GenerateQuery()

	query, _ := boolQuery.Source()
	debug.PrintJson("query", query, true)
}

// 将参数转换成query 值取反
func TestParamsMapList_GenerateQuery2(t *testing.T) {
	params := map[string]interface{}{
		"operate_config_id":                     []interface{}{SignMustNot, 1654},
		"small_course_order.small_course_id":    7094,
		"small_course_order.small_course_stage": 1,
	}

	boolQuery := new(ParamsMapList).LoadConfig("", params).GenerateQuery()

	query, _ := boolQuery.Source()
	debug.PrintJson("query", query, true)
}

// 将参数转换成query 值不存在
func TestParamsMapList_GenerateQuery3(t *testing.T) {
	params := map[string]interface{}{
		"operate_config_id":        1654,
		"small_course_order":       []interface{}{SignMustNot, SignExist}, // 字段不存在
		"intention.intention_type": []interface{}{1, 10},                  // 字段范围是 1<=x<=10
		"customer_id":              []interface{}{1, 10, 123},
	}

	boolQuery := new(ParamsMapList).LoadConfig("", params).GenerateQuery()

	query, _ := boolQuery.Source()
	debug.PrintJson("query", query, true)
}
