package calculate

import (
	"github.com/astaxie/beego/logs"
	"github.com/gaopengfei123123/tools/convert"
	"github.com/olivere/elastic/v7"
	"testing"
)

const MetricsLargeOrder = "MetricsLargeOrderCnt"
const MetricsJoinedClass = "MetricsJoinedClassRoom"

//func getEsCline() *elastic.Client {
//	client, _ := elastic.NewClient(elastic.SetURL("http://0.0.0.0:9200"), elastic.SetTraceLog(new(tracelog)))
//	return client
//}

// 调用简单指标的示例
func TestGetBasicMetrics(t *testing.T) {
	initConfig()
	params := map[string]interface{}{
		"large_course_id":    2138,
		"large_course_stage": 28,
	}

	metrics := []string{
		MetricsLargeOrder,
		MetricsJoinedClass,
	}

	res, err := GetBasicMetrics(nil, metrics, params, getEsCline())
	logs.Info("err: %v", err)

	b, _ := convert.JSONEncode(res)
	logs.Info("res: \n%s", b)
}

func TestGetBasicMetricsWithQuery(t *testing.T) {
	initConfig()
	query := elastic.NewBoolQuery()
	query.Must(elastic.NewTermQuery("large_course_id", 2138))
	query.Must(elastic.NewTermQuery("large_course_stage", 28))

	metrics := []string{
		MetricsLargeOrder,
		MetricsJoinedClass,
	}

	res, err := GetBasicMetricsWithQuery(nil, metrics, query, getEsCline())
	logs.Info("err: %v", err)

	b, _ := convert.JSONEncode(res)
	logs.Info("res: \n%s", b)
}

func TestGetTermsMetrics(t *testing.T) {
	initConfig()
	// 筛选参数
	params := map[string]interface{}{
		"large_course_id":    2138,
		"large_course_stage": 28,
	}

	// 要聚合的指标层级
	termsList := []string{
		"intention.source_key_4",
		"intention.source_key_6",
	}

	// 指标名
	metrics := []string{
		MetricsLargeOrder,
		MetricsJoinedClass,
	}

	res, err := GetTermsMetrics(nil, termsList, metrics, params, getEsCline())
	logs.Info("err: %v", err)

	b, _ := convert.JSONEncode(res)
	logs.Info("res: \n%s", b)
}

// 聚合指标查询1 一级字段+子文档字段
func TestGetTermsMetricsWithQuery(t *testing.T) {
	// 筛选参数
	query := elastic.NewBoolQuery()
	query.Must(elastic.NewTermQuery("large_course_id", 2138))
	query.Must(elastic.NewTermQuery("large_course_stage", 28))

	// 要聚合的指标层级
	termsList := []string{
		"intention_type",
		"intention.source_system",
	}

	// 指标名
	metrics := []string{
		MetricsLargeOrder,
		MetricsJoinedClass,
	}

	res, err := GetTermsMetricsWithQuery(nil, termsList, metrics, query, getEsCline())
	logs.Info("err: %v", err)

	b, _ := convert.JSONEncode(res)
	logs.Info("res: \n%s", b)
}

// 提前注入指标配置
func initConfig() {
	metrics := map[string]AggFunc{
		MetricsLargeOrder:  MetricsLargeOrderCnt,
		MetricsJoinedClass: MetricsJoinedClassRoom,
	}

	for i := range metrics {
		SetMetricsAgg(i, metrics[i])
	}

	SetEsIndex("scrm_clue_new")
}

// MetricsLargeOrderCnt 测试用指标
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

func MetricsJoinedClassRoom(currentTerm ...string) elastic.Aggregation {
	fieldPath := "join_classes"
	field := "join_classes.is_display_group"

	// 指标筛选条件
	termQuery := elastic.NewBoolQuery()
	termQuery.Must(elastic.NewTermQuery("join_classes.is_display_group", 1))
	termQuery.Must(elastic.NewTermsQuery("join_classes.cust_wx_status", 1, 2))

	// 指标聚合方式
	metrics := elastic.NewValueCountAggregation().Field(field)

	return commonNestedMetricsReturn(termQuery, metrics, fieldPath)
}

// 通用的带 query 查询的 nested 返回, 这个函数仅适用于 nested类型聚合 && 存在筛选条件  && 条件是同级子文档中的字段
func commonNestedMetricsReturn(termQuery *elastic.BoolQuery, metrics elastic.Aggregation, fieldPath string) elastic.Aggregation {
	oneFilter := elastic.NewFilterAggregation().Filter(termQuery)
	oneFilter.SubAggregation(SignSingle, metrics)
	metricsNested := elastic.NewNestedAggregation().Path(fieldPath)
	metricsNested.SubAggregation(SignFilter, oneFilter)
	return metricsNested
}
