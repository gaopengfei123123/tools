package calculate

import (
	"github.com/astaxie/beego/logs"
	"github.com/gaopengfei123123/tools/convert"
	"github.com/olivere/elastic/v7"
	"testing"
)

// 测试单指标
const MetricsLargeOrder = "MetricsLargeOrderCnt"
const MetricsJoinedClass = "MetricsJoinedClassRoom"

// 测试直方图聚合
const ClueHistogramLargePayedTimeDateHistogram string = "LargePayedTimeDateHistogram"

func getEsCline() *elastic.Client {
	client, _ := elastic.NewClient(elastic.SetURL("http://0.0.0.0:9200"), elastic.SetTraceLog(new(tracelog)))
	return client
}

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

func TestGetTermsMetrics3(t *testing.T) {
	termQuery := elastic.NewBoolQuery()

	termQuery.Must(elastic.NewTermQuery("form_is_repeat_buy", 1))
	termQuery.Must(elastic.NewTermQuery("operate_config_id", 905))

	ctx, _ := termQuery.Source()
	b, _ := convert.JSONEncode(ctx)
	logs.Info("%s", b)
}

func TestGetBasicMetrics2(t *testing.T) {
	initConfig()
	params := map[string]interface{}{
		"account_id": "2621,2334,14",
		"user_id": []interface{}{ // 范围查询示例
			"0", "888",
		},
		"user_id2": []interface{}{ // 如果不是恰好2个, 就会认为是 terms
			"0", "888", "999",
		},
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

	res, err := GetBasicMetricsWithQuery(nil, nil, metrics, query, getEsCline())
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
		//"large_course_id":    []interface{}{nil, 2138},  // 范围查询方式
	}

	// 要聚合的指标层级
	termsList := []string{
		"intention.source_key_4",
		"intention.source_key_6",
	}

	// 指标名
	metrics := []string{
		//MetricsLargeOrder,
		//MetricsJoinedClass,
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

	res, err := GetTermsMetricsWithQuery(nil, nil, termsList, metrics, query, getEsCline())
	logs.Info("err: %v", err)

	b, _ := convert.JSONEncode(res)
	logs.Info("res: \n%s", b)
}

// 测试获取直方图示例
func TestGetHistogramMetrics(t *testing.T) {
	initConfig()
	// 直方图名
	hisName := ClueHistogramLargePayedTimeDateHistogram
	// 筛选参数
	params := map[string]interface{}{
		"large_course_id":    2138,
		"large_course_stage": 28,
		//"large_pay_time":     [2]int64{1, 123}, // 1 <= x <= 123
	}
	//[2]int{0, 123}, //  x <= 123
	//[2]int{1, 0}, //  1 <= x
	// 指标名
	metrics := []string{
		MetricsLargeOrder,
		MetricsJoinedClass,
	}

	res, err := GetHistogramMetrics(nil, hisName, metrics, params, getEsCline())
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

	histogram := map[string]AggHistogramFunc{
		ClueHistogramLargePayedTimeDateHistogram: HistogramLargePayedTimeDate,
	}

	for i := range histogram {
		SetHistogramAgg(i, histogram[i])
	}

	SetEsIndex("scrm_clue_new")
}

// MetricsLargeOrderCnt 测试用指标
func MetricsLargeOrderCnt(params map[string]interface{}, currentTerm ...string) elastic.Aggregation {
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

func MetricsJoinedClassRoom(params map[string]interface{}, currentTerm ...string) elastic.Aggregation {
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

// HistogramLargePayedTimeDate 直方图统计指标
func HistogramLargePayedTimeDate(aggList map[string]elastic.Aggregation) elastic.Aggregation {
	dataField := "large_pay_time"
	histogram := elastic.NewDateHistogramAggregation().
		Field(dataField).
		Interval("hour").
		Format("YYYY-MM-dd").
		MinDocCount(0).
		TimeZone("Asia/Shanghai")

	// 如果存在指标注入, 就放里面
	if len(aggList) > 0 {
		for aggName, aggBody := range aggList {
			histogram.SubAggregation(aggName, aggBody)
		}
	}
	return histogram
}

// builder 解析器
func TestEsQueryBuilder_LoadParams(t *testing.T) {
	termsList := []string{
		"province_",
	}
	metricsList := []string{}

	params := map[string]interface{}{
		"lib_lib":             "js",
		"corp_type>corp_type": "123",
	}
	builder := new(EsQueryBuilder)
	// 将参数解析成 es query
	requestQuery := builder.LoadParams(termsList, metricsList, params).GetStringQuery()
	logs.Info("requestQuery: %v", requestQuery)
}

// builder 解析器
func TestEsQueryBuilder_LoadParams2(t *testing.T) {
	termsList := []string{
		"properties>account_type", // 针对query 中出现的  properties.account_type 查询
	}
	metricsList := []string{}

	params := map[string]interface{}{
		"event": "siteYsPageView",
	}
	builder := new(EsQueryBuilder)
	// 将参数解析成 es query
	requestQuery := builder.LoadParams(termsList, metricsList, params).GetStringQuery()
	logs.Info("requestQuery: %v", requestQuery)
}

// 解析query
func TestEsQueryBuilder_ParseQuery(t *testing.T) {
	// 筛选参数
	params := map[string]interface{}{
		"large_course_id":    2138,
		"large_course_stage": 28,
		//"large_pay_time":     [2]int64{1, 123}, // 1 <= x <= 123
	}

	// 初始化 builder
	builder := new(EsQueryBuilder)

	query := builder.ParseQuery(params)
	result, _ := convert.JSONEncode(query)
	logs.Info("%s", result)
}

// 解析agg
func TestEsQueryBuilder_ParseAgg(t *testing.T) {
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

	// 初始化 builder
	builder := new(EsQueryBuilder)

	agg := builder.ParseAgg(nil, termsList, metrics)
	result, _ := convert.JSONEncode(agg)
	logs.Info("%s", result)
}

func TestEsQueryBuilder_ParseAgg2(t *testing.T) {
	// 要聚合的指标层级
	termsList := []string{
		//"intention_type",
		"intention>source_system",
		"intention.source_system",
	}

	// 指标名
	metrics := []string{
		MetricsLargeOrder,
		MetricsJoinedClass,
	}

	// 初始化 builder
	builder := new(EsQueryBuilder)

	agg := builder.ParseAgg(nil, termsList, metrics)
	result, _ := convert.JSONEncode(agg)
	logs.Info("%s", result)
}

func TestEsQueryBuilder_ParseParamsToQuery(t *testing.T) {
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
	// 筛选参数
	params := map[string]interface{}{
		"large_course_id":    2138,
		"large_course_stage": 28,
		//"large_pay_time":     [2]int64{1, 123}, // 1 <= x <= 123
	}

	// 初始化 builder
	builder := new(EsQueryBuilder)

	query := builder.ParseParamsToQuery(termsList, metrics, params)
	result, _ := convert.JSONEncode(query)
	logs.Info("%s", result)
}

// 模拟params 传空的情况
func TestEsQueryBuilder_ParseParamsToQuery2(t *testing.T) {
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
	// 筛选参数
	params := map[string]interface{}{
		"large_course_id":    2138,
		"large_course_stage": 28,
		//"large_pay_time":     [2]int64{1, 123}, // 1 <= x <= 123
	}

	// 初始化 builder
	builder := new(EsQueryBuilder)

	query := builder.ParseParamsToQuery(termsList, metrics, params)
	result, _ := convert.JSONEncode(query)
	logs.Info("%s", result)
}

// builder 解析器
func TestEsQueryBuilder_ParseSearchResult(t *testing.T) {
	termsList := []string{
		"province_",
	}
	metricsList := []string{}
	params := map[string]interface{}{
		"lib_lib": "js",
	}

	// 初始化 builder
	builder := new(EsQueryBuilder)
	builder.LoadParams(termsList, metricsList, params)

	// 不管任何途径获取道德 es 结果, 只要符合 *elastic.SearchResult 就行
	resultJson := `{"_shards":{"total":3,"failed":0,"successful":3,"skipped":0},"hits":{"hits":[],"total":156,"max_score":0},"took":0,"timed_out":false,"aggregations":{"term_province_":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"doc_count":156,"key":"北京"}]}}}`
	searchResult := &elastic.SearchResult{}
	convert.JSONDecode([]byte(resultJson), searchResult)

	// 解析搜索结果
	result, err := builder.ParseSearchResult(searchResult)
	jsb, _ := convert.JSONEncode(result)
	logs.Info("result: %s, err: %v", jsb, err)
}

func TestEsQueryBuilder_GetStringQuery(t *testing.T) {
	metricsList := []string{
		MetricsLargeOrder,
		MetricsJoinedClass,
	}

	params := map[string]interface{}{
		"event": "siteYsPageView",
	}
	builder := new(EsQueryBuilder)
	tmp := builder.LoadParams(nil, metricsList, params)

	ss := tmp.GetStringQuery()
	logs.Info("ss: %s", ss)
}
