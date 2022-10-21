package calculate

import (
	"github.com/astaxie/beego/logs"
	"github.com/gaopengfei123123/tools/convert"
	"github.com/olivere/elastic/v7"
	"testing"
)

const MetricsLargeOrder = "MetricsLargeOrderCnt"

// 调用简单指标的示例
func TestGetBasicMetrics(t *testing.T) {
	initConfig()
	params := map[string]interface{}{
		"large_course_id":    2138,
		"large_course_stage": 28,
	}

	metrics := []string{
		"MetricsLargeOrderCnt",
	}

	res, err := GetBasicMetrics(metrics, params, getEsCline(), nil)
	logs.Info("err: %v", err)

	b, _ := convert.JSONEncode(res)
	logs.Info("res: \n%s", b)
}

func TestGetBasicMetricsWithQuery(t *testing.T) {
	query := elastic.NewBoolQuery()
	query.Must(elastic.NewTermQuery("large_course_id", 2138))
	query.Must(elastic.NewTermQuery("large_course_stage", 28))

	metrics := []string{
		MetricsLargeOrder,
	}

	res, err := GetBasicMetricsWithQuery(metrics, query, getEsCline(), nil)
	logs.Info("err: %v", err)

	b, _ := convert.JSONEncode(res)
	logs.Info("res: \n%s", b)
}

func initConfig() {
	metrics := map[string]AggFunc{
		MetricsLargeOrder: MetricsLargeOrderCnt,
	}

	for i := range metrics {
		SetMetricsAgg(i, metrics[i])
	}

	SetEsIndex("scrm_clue_new")
}

// MetricsLargeOrderCnt
func MetricsLargeOrderCnt(currentTerm ...string) elastic.Aggregation {
	termQuery := elastic.NewBoolQuery()

	// 筛选条件
	termQuery.Must(elastic.NewTermQuery("is_buy_large_course", 1))
	termQuery.Must(elastic.NewTermQuery("is_submit_large_order", 1))
	metrics := elastic.NewFilterAggregation().Filter(termQuery)

	// 聚合方式以及字段
	aggCount := elastic.NewValueCountAggregation().Field("is_buy_large_course")
	metrics.SubAggregation(SignSingle, aggCount)
	return metrics
}
