package debug

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/olivere/elastic/v7"
)

// PrintEsBoolQueryAgg 打印聚合函数
func PrintEsBoolQueryAgg(query *elastic.BoolQuery, aggs ...elastic.Aggregation) []byte {
	qsrc, _ := query.Source()

	// 带 agg 的
	if len(aggs) > 0 && aggs[0] != nil {
		agg := aggs[0]
		asrc, _ := agg.Source()
		tmp := map[string]interface{}{
			"query": qsrc,
			"aggregations": map[string]interface{}{
				"data": asrc,
			},
			"size": 0,
		}
		tmpB, _ := json.Marshal(tmp)
		logs.Debug("PrintEsBoolQuery: %s", tmpB)
		return tmpB
	}

	// 只有 query 的
	tmp := map[string]interface{}{
		"query": qsrc,
		"size":  0,
	}
	tmpB, _ := json.Marshal(tmp)
	logs.Debug("PrintEsBoolQuery: %s", tmpB)
	return tmpB
}

func PrintAgg(aggs ...elastic.Aggregation) []byte {
	result := make(map[string]interface{})

	for k, v := range aggs {
		key := fmt.Sprintf("agg_%v", k)
		src, err := v.Source()
		if err != nil {
			continue
		}
		result[key] = src
	}
	tmpB, _ := json.Marshal(result)
	logs.Debug("PrintEsBoolQuery: %s", tmpB)
	return tmpB
}
