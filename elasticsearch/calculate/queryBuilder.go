package calculate

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/olivere/elastic/v7"
)

// 这里把参数转 es query 和搜索结果转返回值抽象出来

type EsQueryBuilder struct {
	RequestQuery interface{} `json:"request_query"`
	termsList    []string
	metricsList  []string
	params       map[string]interface{}
}

// ParseSearchResult 解析结果, 得是从 builder 生成出去的 query 返回的结果这里能解析出来
func (builder *EsQueryBuilder) ParseSearchResult(searchResult *elastic.SearchResult) (result TermResult, err error) {
	if searchResult.Error != nil && searchResult.Error.Reason != "" {
		err = fmt.Errorf(searchResult.Error.Reason)
		return
	}

	return BatchGetValueFromTerms(searchResult.Aggregations, builder.metricsList, builder.termsList)
}

// LoadParams 这里是一切的入口
func (builder *EsQueryBuilder) LoadParams(termsList []string, metricsList []string, params map[string]interface{}) *EsQueryBuilder {
	builder.termsList = termsList
	builder.metricsList = metricsList
	builder.params = params
	builder.RequestQuery = builder.ParseParamsToQuery(termsList, metricsList, params)
	return builder
}

func (builder *EsQueryBuilder) GetStringQuery() string {
	if builder.RequestQuery == nil {
		return ""
	}
	b, _ := json.Marshal(builder.RequestQuery)
	return string(b)
}

// ParseParamsToQuery 返回完整的 terms 聚合query, 直接用到es查询中
func (builder *EsQueryBuilder) ParseParamsToQuery(termsList []string, metricsList []string, params map[string]interface{}) interface{} {
	agg := builder.ParseAgg(termsList, metricsList)
	logs.Trace("agg:%v", agg)

	query := builder.ParseQuery(params)
	logs.Trace("query: %v", query)

	result := map[string]interface{}{}

	if query != nil {
		result["query"] = query
	}
	if agg != nil {
		result["aggregations"] = agg
	}
	result["size"] = 0

	builder.RequestQuery = result

	return result
}

// ParseQuery 解析筛选主文档的 query
func (builder *EsQueryBuilder) ParseQuery(params map[string]interface{}) interface{} {
	query := new(ParamsMapList).LoadConfig("", params).GenerateQuery()
	queryInterface, _ := query.Source()
	return queryInterface
}

// ParseAgg 解析 agg
func (builder *EsQueryBuilder) ParseAgg(termsList []string, metricsList []string) interface{} {
	agg := BuildTermAgg("", termsList, metricsList)
	service := elastic.NewTermsAggregation()
	_, key, isNested := CheckKeyNested(termsList[0])
	if isNested {
		service.SubAggregation(SignChild, agg)
	} else {
		service.SubAggregation(GenTermKey(key), agg)
	}
	serviceInterface, _ := service.Source()

	aggMap, ok := serviceInterface.(map[string]interface{})
	if !ok {
		return nil
	}
	return aggMap["aggregations"]
}
