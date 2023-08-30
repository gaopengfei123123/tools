package calculate

import (
	"github.com/olivere/elastic/v7"
	"sync"
)

const (
	SignAggList       = "metrics_list"
	SignSingle        = "metric_data"
	SignFilter        = "metric_filter"
	SignFilterReverse = "metric_filter_reverse"
	SignTerms         = "metrics_terms"
	SignChild         = "metrics_child"
	SignNested        = "."              // nested类型字段连接方式  nested A下的字段B, 写成 A.B 实际按
	SignObject        = ">"              // object类型字段链接方式  object A下的字段 B, 写成 A>B, 最终转换成 A.B
	SignTimeRange     = "agg_time_range" // 从外来参数传入到指标中的参数标识
)

const (
	TypeTerms         = "terms"
	TypeCount         = "count"
	TypeTerm          = "term"
	TypeSum           = "sum"
	TypeAvg           = "avg"
	TypeNested        = "nested"
	QueryMust         = "must"
	QueryRange        = "range"
	QueryMustMulti    = "mustMulti"
	QueryMustNot      = "MustNot"
	QueryMustNotMulti = "mustNotMulti"
	SignNotNull       = "NOT_NULL" // 非空标识符

	MaxSize = 10000000 // 聚合 terms 的时候的最大数量
)

func init() {
	esconfig = &ESConfig{}
	esconfig.AggFuncList = make(map[string]AggFunc)
}

type ESConfig struct {
	EsIndex        string                      // 要用到的 es 文档
	AggFuncList    map[string]AggFunc          // 存放 AggFunc 函数
	HistogramList  map[string]AggHistogramFunc // 存放直方图聚合指标
	GetAggListFunc GetAggFunc
	sync.Mutex
}

var esconfig *ESConfig

// 检测是否需要返回根目录
var checkAggNeedReverse CheckAggReverseNest

// AggFunc 指标查询函数要求的格式
type AggFunc func(params map[string]interface{}, currentTerm ...string) elastic.Aggregation
type GetAggFunc func(metricName string, sceneName ...string) AggFunc       // 外部注入的获取指标的方法, 如果没有的话, 就默认读取AggFuncList这里的指标
type CheckAggReverseNest func(metricName string, sceneName ...string) bool // 检测函数是否需要反转 true的时候, 不执行 reverse_nest 操作

// AggHistogramFunc 直方聚合查询
type AggHistogramFunc func(aggList map[string]elastic.Aggregation) elastic.Aggregation

func GetEsConfig() *ESConfig {
	return esconfig
}

func SetMetricsAgg(metricName string, aggFuncList ...AggFunc) {
	esconfig.Lock()
	defer esconfig.Unlock()

	if esconfig.AggFuncList == nil {
		esconfig.AggFuncList = make(map[string]AggFunc)
	}

	for i := range aggFuncList {
		esconfig.AggFuncList[metricName] = aggFuncList[i]
	}
}

// SetMetricsAggFunc 设置获取指标信息的方法
func SetMetricsAggFunc(fn GetAggFunc) {
	esconfig.GetAggListFunc = fn
}

func SetMetricsReverse(fn CheckAggReverseNest) {
	checkAggNeedReverse = fn
}

// CheckAggNeeReverse 检测该指标是否就是需要在 nested 环境下使用 true 则不执行 reverse_nested
func CheckAggNeeReverse(metricName string, sceneName ...string) bool {
	if checkAggNeedReverse == nil {
		return false
	}

	return checkAggNeedReverse(metricName, sceneName...)
}

// GetMetricsAgg 获取注入的指标
func (ec *ESConfig) GetMetricsAgg(metricName string, sceneName ...string) AggFunc {
	// 优先使用函数方法
	if ec.GetAggListFunc != nil {
		agg := ec.GetAggListFunc(metricName, sceneName...)
		if agg != nil {
			return agg
		}
	}

	ec.Lock()
	defer ec.Unlock()
	aggFunc, ok := ec.AggFuncList[metricName]
	if ok {
		return aggFunc
	}
	return nil
}

// SetEsIndex 确定要查询的esindex
func SetEsIndex(index string) {
	esconfig.EsIndex = index
}

func (ec *ESConfig) GetEsIndex(params map[string]interface{}, scene ...string) string {
	return ec.EsIndex
}

// SetHistogramAgg 设置
func SetHistogramAgg(histogramName string, metricsFunc AggHistogramFunc) {
	esconfig.Lock()
	defer esconfig.Unlock()

	if esconfig.HistogramList == nil {
		esconfig.HistogramList = make(map[string]AggHistogramFunc)
	}
	esconfig.HistogramList[histogramName] = metricsFunc
	return
}

func (ec *ESConfig) GetHistogramAgg(histogramName string, metricsList map[string]elastic.Aggregation, sceneName ...string) elastic.Aggregation {
	ec.Lock()
	defer ec.Unlock()

	aggFunc, ok := ec.HistogramList[histogramName]
	if !ok {
		return nil
	}
	return aggFunc(metricsList)
}
