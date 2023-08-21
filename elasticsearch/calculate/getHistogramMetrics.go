package calculate

import (
	"context"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/gaopengfei123123/tools/convert"
	"github.com/olivere/elastic/v7"
)

type HistogramResult struct {
	HistogramName string
	MetricsList   []string
	Result        []HistogramItem
}
type HistogramItem struct {
	KeyAsString string
	Key         interface{}
	DocCount    int64
	Metrics     map[string]interface{}
}

// LoadDataFromAggregations 从返回结果中解析直方图的结果
func (hr *HistogramResult) LoadDataFromAggregations(aggRes elastic.Aggregations) (HistogramResult, error) {
	bucketList, exist := aggRes.Histogram(hr.HistogramName)

	if !exist {
		return *hr, fmt.Errorf("搜索结果找不到对应的 histogram: %s", hr.HistogramName)
	}

	hr.Result = make([]HistogramItem, 0, len(bucketList.Buckets))
	for index := range bucketList.Buckets {
		cur := bucketList.Buckets[index]
		tmp := HistogramItem{}
		err := tmp.LoadDataFromItem(cur, hr.MetricsList)
		if err != nil {
			continue
		}
		hr.Result = append(hr.Result, tmp)
	}
	return *hr, nil
}

func (hi *HistogramItem) LoadDataFromItem(bucketItem *elastic.AggregationBucketHistogramItem, metricsList []string) error {
	hi.Key = bucketItem.Key
	hi.KeyAsString = *bucketItem.KeyAsString
	hi.DocCount = bucketItem.DocCount

	if bucketItem.Aggregations == nil || len(metricsList) == 0 {
		return nil
	}
	// 存在指标的时候读指标
	dataList, err := BatchGetValueFromAggData(bucketItem.Aggregations, metricsList)
	if err != nil {
		return err
	}
	hi.Metrics = dataList
	return nil
}

// GetHistogramMetrics 获取直方图聚合指标
func GetHistogramMetrics(ctx context.Context, histogramName string, metricsList []string, params map[string]interface{}, client *elastic.Client, sceneName ...string) (result HistogramResult, err error) {
	var scene string
	if len(sceneName) != 0 {
		scene = sceneName[0]
	}
	query := new(ParamsMapList).LoadConfig(scene, params).GenerateQuery()
	return GetHistogramMetricsWithQuery(ctx, histogramName, params, metricsList, query, client, sceneName...)
}

// GetHistogramMetricsWithQuery 获取直方图, 参数由 query 格式传入
func GetHistogramMetricsWithQuery(ctx context.Context, histogramName string, params map[string]interface{}, metricsList []string, query elastic.Query, client *elastic.Client, sceneName ...string) (result HistogramResult, err error) {
	if client == nil {
		err = fmt.Errorf("esClient is nil")
		return
	}

	if len(metricsList) == 0 {
		err = fmt.Errorf("metricsList is nil")
		return
	}

	if ctx == nil {
		ctx = context.Background()
	}

	var scene string
	if len(sceneName) != 0 {
		scene = sceneName[0]
	}
	logs.Info("scene: %v", scene)

	result = HistogramResult{
		HistogramName: histogramName,
		MetricsList:   metricsList,
	}

	histogramAgg := BuildHistogramMetricsAgg(histogramName, params, metricsList, scene)

	// 获取需要用到的索引名
	esIndex := esconfig.GetEsIndex(params, scene)
	service := client.Search().Index(esIndex)
	service.Aggregation(histogramName, histogramAgg)
	service.Query(query)

	searchResult, err := service.Size(0).Do(ctx)

	b, _ := convert.JSONEncode(searchResult)
	logs.Info("searchResult: %s, err: %v", b, err)

	if searchResult == nil || searchResult.Aggregations == nil {
		return result, err
	}

	return result.LoadDataFromAggregations(searchResult.Aggregations)
}

// BuildHistogramMetricsAgg 构建直方图的agg
func BuildHistogramMetricsAgg(histogramName string, params map[string]interface{}, metricsList []string, sceneName ...string) elastic.Aggregation {
	var scene string
	if len(sceneName) != 0 {
		scene = sceneName[0]
	}

	metricsBodyList := make(map[string]elastic.Aggregation)
	// 初始化聚合指标
	for i := range metricsList {
		key := metricsList[i]
		aggFunc := esconfig.GetMetricsAgg(key, scene)
		metricsBodyList[key] = aggFunc(params)
	}

	return esconfig.GetHistogramAgg(histogramName, metricsBodyList)
}
