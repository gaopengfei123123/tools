package calculate

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"strings"
)

type TermResult struct {
	Result      []TermItem
	TermsList   []string
	MetricsList []string
}
type TermItem struct {
	Terms []interface{}
	Cnt   int64
	Data  map[string]interface{}
	agg   *elastic.Aggregations // 不对外暴露
}

// GetFilteredData 将 nil 值转成0
func (t *TermItem) GetFilteredData() map[string]interface{} {
	for k, v := range t.Data {
		if v == nil {
			t.Data[k] = 0
		}
	}
	return t.Data
}

// GetCombineTerm 获取
func (t *TermItem) GetCombineTerm() string {
	str := make([]string, 0)

	for i := 0; i < len(t.Terms); i++ {
		str = append(str, fmt.Sprintf("%v", t.Terms[i]))
	}
	return strings.Join(str, "_")
}

func (ti *TermResult) GetMapDataResult() (result []map[string]interface{}, err error) {
	result = make([]map[string]interface{}, 0)
	if len(ti.Result) == 0 {
		return nil, nil
	}

	for _, curTermItem := range ti.Result {
		term := curTermItem.GetCombineTerm()
		metrics := curTermItem.GetFilteredData()
		tmp := map[string]interface{}{
			"metrics": metrics,
			"terms":   term,
		}
		result = append(result, tmp)
	}

	return result, nil
}

// LoopTermsFromAgg 当传入到这里的时候, 已经能取一级 term 的 bucket 了
func (ti *TermResult) LoopTermsFromAgg(aggRes *elastic.AggregationBucketKeyItems, columnValue []interface{}, level ...int) error {
	var lv int // 这里标明当前所处的层级, 初始是1
	if len(level) > 0 {
		lv = level[0]
	} else {
		lv = 1
	}

	// 初始化
	if lv == 1 {
		columnValue = make([]interface{}, 0, len(ti.TermsList))
	}

	childTerm, exist := aggRes.Terms(SignChild)
	// nested 类型的, 将从子agg 中提取
	if exist {
		aggRes = childTerm
	}

	for _, bucket := range aggRes.Buckets {
		k := bucket.Key
		// 新生成一个字段数组
		tmpColumnValue := append(columnValue, k)

		// 这时说明已经遍历到最下层
		if lv == len(ti.TermsList) {
			tmpColumn := make([]interface{}, len(tmpColumnValue))
			copy(tmpColumn, tmpColumnValue) // 切片深拷贝, 避免内存共享的干扰 -- 一个 go 语言的暗坑

			metricsData, _ := BatchGetValueFromAggData(bucket.Aggregations, ti.MetricsList)
			//logs.Info("metricsData: %v, err: %v", metricsData, err)
			tmp := TermItem{
				Terms: tmpColumn,
				Data:  metricsData,
				Cnt:   bucket.DocCount,
			}
			ti.Result = append(ti.Result, tmp)
			continue
		}

		// 判断当前这个 key 是不是nested 类型
		_, key, isNested := checkKeyNested(ti.TermsList[lv])
		var tmpAgg *elastic.AggregationBucketKeyItems
		if isNested {
			tmpAgg, exist = bucket.Terms(SignChild)
			if exist {
				curKey := genTermKey(key)
				tmpAgg, exist = tmpAgg.Terms(curKey)
			} else {
				curKey := genTermKey(key)
				tmpAgg, exist = bucket.Terms(curKey)
			}
			//logs.Info("tmpAgg: %s", tmpAgg)
		} else {
			curKey := genTermKey(key)

			tmpAgg, exist = bucket.Aggregations.Terms(curKey)

			// 这里如果存在子页面嵌套的情况, 就多读一层, 测试用例 TestGetTermsMetrics5
			ctmpAgg, cexist := tmpAgg.Aggregations.Terms(curKey)
			if cexist {
				tmpAgg = ctmpAgg
			}
		}

		if exist {
			err := ti.LoopTermsFromAgg(tmpAgg, tmpColumnValue, lv+1)
			if err != nil {
				return nil
			}
		}
	}

	return nil
}

// GetTermsMetrics 获取多级指标
func GetTermsMetrics(ctx context.Context, termsList []string, metricsList []string, params map[string]interface{}, client *elastic.Client, sceneName ...string) (result TermResult, err error) {
	var scene string
	if len(sceneName) != 0 {
		scene = sceneName[0]
	}
	// 获取 query
	query := new(ParamsMapList).LoadConfig(scene, params).GenerateQuery()
	return GetTermsMetricsWithQuery(ctx, termsList, metricsList, query, client, sceneName...)
}

// GetTermsMetricsWithQuery 以原生传参的方式, 获取聚合指标
func GetTermsMetricsWithQuery(ctx context.Context, termsList []string, metricsList []string, query elastic.Query, client *elastic.Client, sceneName ...string) (result TermResult, err error) {
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

	result = TermResult{}

	// 获取需要用到的索引名
	esIndex := esconfig.GetEsIndex(scene)
	service := client.Search().Index(esIndex)

	// 递归组装 term 以及 agg
	agg := BuildTermAgg(scene, termsList, metricsList)

	if len(termsList) == 0 {
		return result, fmt.Errorf("请输入要聚合的字段")
	}

	_, key, isNested := checkKeyNested(termsList[0])
	if isNested {
		service.Aggregation(SignChild, agg)
	} else {
		service.Aggregation(genTermKey(key), agg)
	}

	service.Query(query)

	searchResult, err := service.Size(0).Do(ctx)

	if searchResult == nil || searchResult.Aggregations == nil {
		return result, err
	}

	return BatchGetValueFromTerms(searchResult.Aggregations, metricsList, termsList)
}

// BatchGenerateMetricsAgg 将指标批量加载进去
func BatchGenerateMetricsAgg(scene string, agg *elastic.TermsAggregation, metricsList []string, currentTerm string, isNested bool) *elastic.TermsAggregation {
	if agg == nil || len(metricsList) == 0 {
		return agg
	}

	//_, _, isNested := checkKeyNested(currentTerm)

	for i := 0; i < len(metricsList); i++ {
		metricsName := metricsList[i]
		aggFunc := esconfig.GetMetricsAgg(metricsName, scene)
		if aggFunc == nil {
			continue
		}

		// 如果是 nested 类型数据, 需要回退一格
		if isNested {
			// 回退一格到主文档
			reverseAgg := elastic.NewReverseNestedAggregation()
			reverseAgg.SubAggregation(metricsName, aggFunc())
			agg.SubAggregation(metricsName, reverseAgg)
		} else {
			agg.SubAggregation(metricsName, aggFunc())
		}
	}
	return agg
}

// BuildTermAgg 递归生成多级 term
func BuildTermAgg(scene string, termList []string, metricsList []string, level ...int) elastic.Aggregation {
	var lv int
	if len(level) != 0 {
		lv = level[0]
	} else {
		lv = 0
	}
	// 终止递归
	if len(termList) == 0 || len(termList) == lv {
		return nil
	}

	pth, key, isNested := checkKeyNested(termList[lv])

	agg := elastic.NewTermsAggregation().Field(key).Size(MaxSize).OrderByKeyAsc()

	// 这里递归生成嵌套的 term
	lv = lv + 1
	// 这时候说明是最后一层term 了, 给它挂上指标
	if lv == len(termList) {
		agg = BatchGenerateMetricsAgg(scene, agg, metricsList, key, isNested)
	}

	childAgg := BuildTermAgg(scene, termList, metricsList, lv)
	if childAgg != nil {
		// 这里是父级term, 需要看下级是否和自己path 一致, 一致就不生成 esconfig.SignChild
		if lv < len(termList)+1 {
			cpth, ckey, cIsNested := checkKeyNested(termList[lv])
			// 如果上级 path 和当前想通, 则不加 nested
			if cpth == pth {
				agg.SubAggregation(genTermKey(ckey), childAgg)
			} else {
				if cIsNested {
					agg.SubAggregation(SignChild, childAgg)
				} else {
					agg.SubAggregation(genTermKey(ckey), childAgg)
				}
			}
		}
	}

	// 这里是当前term
	if isNested {
		// 这里从1开始是因为前面已经+1了, 需要排除掉第一个 nested 类型
		if lv > 1 {
			ppth, _, parentIsNested := checkKeyNested(termList[lv-2])
			//ppth, _, parentIsNested := checkKeyNested(termList[lv-2])
			// 父类如果是nested, 且不是相同子文档下的, 那么当前的也得用 reverse_nested
			if parentIsNested && ppth != pth {
				aggNested := elastic.NewReverseNestedAggregation().Path(pth)
				aggNested.SubAggregation(genTermKey(key), agg)
				return aggNested
			}

			// 下级 term 和上级属于同一个 term, 不用加 nested
			if ppth == pth {
				return agg
			}
		}
		aggNested := elastic.NewNestedAggregation().Path(pth)
		aggNested.SubAggregation(genTermKey(key), agg)
		return aggNested
	}

	// 针对二级以上的 term, 当前是一级字段, 父级是个子文档
	if lv > 1 {
		ppth, _, _ := checkKeyNested(termList[lv-2])
		if ppth != pth {
			aggNested := elastic.NewReverseNestedAggregation()
			aggNested.SubAggregation(genTermKey(key), agg)
			return aggNested
		}
	}

	return agg
}

// BatchGetValueFromTerms 从聚合 terms 中批量获取值
func BatchGetValueFromTerms(aggRes elastic.Aggregations, metricList []string, termsList []string) (result TermResult, err error) {
	tmp := make([]TermItem, 0)
	result.Result = tmp
	if aggRes == nil {
		err = fmt.Errorf("agg 结果为空")
		return
	}

	if len(termsList) == 0 {
		err = fmt.Errorf("terms 为空")
		return
	}

	// 因为nested 类型数据上层始终有一层  esconfig.SignChild 在, 所以取第一个值得方法也不同
	_, key, isNested := checkKeyNested(termsList[0])
	var agg *elastic.AggregationBucketKeyItems
	var exist bool
	if isNested {
		agg, exist = aggRes.Terms(SignChild)
		if !exist {
			err = fmt.Errorf("agg 找不到 nested 数据")
			return
		}
		agg, exist = agg.Terms(genTermKey(key))
	} else {
		agg, exist = aggRes.Terms(genTermKey(key))
	}

	if !exist {
		err = fmt.Errorf("agg 找不到")
		return
	}
	result = TermResult{
		Result:      make([]TermItem, 0),
		TermsList:   termsList,
		MetricsList: metricList,
	}
	// 第一步传空, 后续递归叠加
	err = result.LoopTermsFromAgg(agg, nil)
	//logs.Info("err: %v", err)
	//logs.Info("BatchGetValueFromTerms: %s", agg)
	//
	//b, _ := convert.JSONEncode(res)
	//logs.Info("result: %s", b)
	return result, err
}

// 检测字段是否属于子文档中的字段
func checkKeyNested(key string) (pth string, ky string, isNested bool) {
	ky = key
	isNested = false
	arr := strings.Split(key, SignNested)
	if len(arr) > 1 {
		isNested = true
		pth = arr[0]
		return
	}

	// 这里针对非 nested 类型的 object 类型数据
	arr = strings.Split(key, SignObject)
	if len(arr) > 1 {
		isNested = false
		ky = strings.Join(arr, ".")
		pth = ""
	}
	return
}

// 组合 term聚合的
func genTermKey(key string) string {
	return "term_" + key
}

// GenTermKey 组合 term聚合的 key
func GenTermKey(key string) string {
	return genTermKey(key)
}

func CheckKeyNested(key string) (pth string, ky string, isNested bool) {
	return checkKeyNested(key)
}
