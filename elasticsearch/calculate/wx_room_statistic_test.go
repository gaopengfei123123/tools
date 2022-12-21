package calculate

import (
	"github.com/olivere/elastic/v7"
	"testing"
)

const MetricsRoomOrder = "MetricsRoomOrder"

// 提前注入指标配置
func initWxRoomConfig() {
	metrics := map[string]AggFunc{
		MetricsRoomOrder: MetricsRoomOrderFunc,
	}

	for i := range metrics {
		SetMetricsAgg(i, metrics[i])
	}
	SetEsIndex("wx_room_statistic_month")
}

// MetricsLargeOrderCnt 测试用指标
func MetricsRoomOrderFunc(currentTerm ...string) elastic.Aggregation {
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

func TestGetTermsMetrics2(t *testing.T) {

}
