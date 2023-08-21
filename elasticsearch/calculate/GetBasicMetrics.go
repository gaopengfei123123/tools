package calculate

import (
	"context"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/olivere/elastic/v7"
)

/* es 聚合业务代码 */

// GetBasicMetrics 获取简单指标
func GetBasicMetrics(ctx context.Context, metricsList []string, params map[string]interface{}, client *elastic.Client, sceneName ...string) (result map[string]interface{}, err error) {
	var scene string
	if len(sceneName) != 0 {
		scene = sceneName[0]
	}

	// 这里获取 query
	query := new(ParamsMapList).LoadConfig(scene, params).GenerateQuery()
	return GetBasicMetricsWithQuery(ctx, params, metricsList, query, client, sceneName...)
}

// GetBasicMetricsWithQuery 以原生传query 的方式获取值
func GetBasicMetricsWithQuery(ctx context.Context, params map[string]interface{}, metricsList []string, query elastic.Query, client *elastic.Client, sceneName ...string) (result map[string]interface{}, err error) {
	if client == nil {
		err = fmt.Errorf("esClient is nil")
		return
	}

	if ctx == nil {
		ctx = context.Background()
	}

	var scene string
	if len(sceneName) != 0 {
		scene = sceneName[0]
	}

	result = make(map[string]interface{})
	if len(metricsList) == 0 {
		return
	}

	esIndex := esconfig.GetEsIndex(params, scene)
	service := client.Search().Index(esIndex)

	// 这里将循环获取指标函数
	for i := 0; i < len(metricsList); i++ {
		name := metricsList[i]
		aggFunc := esconfig.GetMetricsAgg(name, scene)
		if aggFunc == nil {
			continue
		}
		service.Aggregation(name, aggFunc(params))
	}

	service.Query(query)
	searchResult, err := service.Size(0).Do(ctx)

	if searchResult == nil || searchResult.Aggregations == nil {
		return result, err
	}

	return BatchGetValueFromAggData(searchResult.Aggregations, metricsList)
}

// BatchGetValueFromAggData 从 agg 中提取数据
func BatchGetValueFromAggData(aggRes elastic.Aggregations, metricList []string) (result map[string]interface{}, err error) {
	if len(aggRes) == 0 {
		err = fmt.Errorf("不存在 Agg 数据")
		return nil, err
	}

	result = make(map[string]interface{})

	for i := 0; i < len(metricList); i++ {
		name := metricList[i]
		v := GetValueFromAggData(aggRes, name)
		result[name] = v
	}

	return result, nil
}

// GetValueFromAggData 单个指标从 agg 中取数
func GetValueFromAggData(aggRes elastic.Aggregations, metricsName string) (result interface{}) {
	bucketKeyItem, exist := aggRes.Terms(metricsName)

	if !exist {
		return nil
	}

	// 如果回退的情况, 就再往里挖一层
	tBucket, tExist := bucketKeyItem.Terms(metricsName)
	if tExist {
		bucketKeyItem = tBucket
	}

	// 检查是否存在metric_filter
	bucketKeyItem2, exist := bucketKeyItem.Terms(SignFilter)
	if exist {
		// 存在filter 的情况下, 将这一层级提出来
		bucketKeyItem = bucketKeyItem2
	}

	// 看 sdk 源码里, valueCount 和 sum 的实现是一样的, 所以这里不做区分, 统一返回 float64
	terms, exist := bucketKeyItem.ValueCount(SignSingle)
	if !exist || terms.Value == nil {
		return nil
	}
	return *terms.Value
}

// 测试环境的 es
type tracelog struct{}

func (tracelog) Printf(format string, v ...interface{}) {
	logs.Trace("elasticPrint:")
	logs.Trace(format, v...)
}
