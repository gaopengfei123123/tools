package calculate

import (
	"github.com/olivere/elastic/v7"
	"sync"
)

const (
	SignAggList = "metrics_list"
	SignSingle  = "metric_data"
	SignFilter  = "metric_filter"
	SignTerms   = "metrics_terms"
	SignChild   = "metrics_child"
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

	MaxSize = 10000000 // 聚合 terms 的时候的最大数量
)

func init() {
	esconfig = &ESConfig{}
	esconfig.AggFuncList = make(map[string]AggFunc)
}

type ESConfig struct {
	EsIndex       string                      // 要用到的 es 文档
	AggFuncList   map[string]AggFunc          // 存放 AggFunc 函数
	HistogramList map[string]AggHistogramFunc // 存放直方图聚合指标
	sync.Mutex
}

var esconfig *ESConfig

// AggFunc 指标查询函数要求的格式
type AggFunc func(currentTerm ...string) elastic.Aggregation

// AggHistogramFunc 直方聚合查询
type AggHistogramFunc func(aggList map[string]elastic.Aggregation) elastic.Aggregation

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

// GetMetricsAgg 获取注入的指标
func (ec *ESConfig) GetMetricsAgg(metricName string, sceneName ...string) AggFunc {
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

func (ec *ESConfig) GetEsIndex(scene ...string) string {
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
