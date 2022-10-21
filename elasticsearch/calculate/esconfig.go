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
	AggFuncList map[string]AggFunc // 存放 AggFunc 函数
	EsIndex     string             // 要用到的 es 文档

	sync.Mutex
}

var esconfig *ESConfig

// AggFunc 指标查询函数要求的格式
type AggFunc func(currentTerm ...string) elastic.Aggregation

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
