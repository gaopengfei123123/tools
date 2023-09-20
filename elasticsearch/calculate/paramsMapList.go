package calculate

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/olivere/elastic/v7"
	"strings"
)

// EsTerm 指标配置, 只允许一个字段值做指标
type EsTerm struct {
	ConfigType string //指标的类型 terms
	ColumnType string //字段类型  term/nested
	Path       string // 字段是 nested 时的路径, 如 join_classes
	ColumnName string // 字段 key 值, 例如 join_classes.cust_wxid
}

type EsAggregation struct {
	ESIndex       string                `json:"es_index"`
	Type          string                `json:"type"`         // 聚合类型  count/terms/avg
	MetricsName   string                `json:"metrics_name"` // 指标名称
	Field         string                `json:"field"`        // 要聚合的字段名
	FieldPath     string                `json:"filed_path"`   // 字段路径
	Size          int                   `json:"size"`         // 数量限制 1000
	FieldType     string                `json:"filed_type"`   // 字段类型 term/nested
	TermQuery     *EsQueryConfig        `json:"term_query"`   // 如果是type是 terms 的, 需要指定筛选条件
	SearchQuery   *EsQueryConfig        `json:"search_query"` // 整体数据的搜索条件
	searchResult  *elastic.SearchResult // es 查询结果
	err           error                 // 返回错误
	ParamsMapList *ParamsMapList        // 初始化后的参数配置
}

type EsQueryConfig struct {
	Conditions []*EsQueryCondition
}

// EsQueryCondition es 搜索条件
type EsQueryCondition struct {
	Type       string // key 类型   term/terms/nested
	Path       string // nested 的路径  join_classes
	Conditions []*SearchCondition
}

type SearchCondition struct {
	Type     string      // term/terms/range/
	Key      string      // 字段名  join_classes.cust_wxid
	Optional int         // 是否可选, 1 可选, 0 必填
	Value    interface{} //字段值  xxxwww
}

// ParamsMapList 参数映射图
type ParamsMapList struct {
	SceneName string `json:"scene_name"`
	//MetricsName   string                      `json:"metrics_name"`
	Params        map[string]string           `json:"params"`
	NestedGroup   map[string]map[string]bool  // 结构是 nested path: []key
	OriginList    map[string]*SearchCondition // 指定的 key => 最终生效的查询
	EsQueryConfig *EsQueryConfig              // 最终生成的 query 组合
}

func (pl *ParamsMapList) LoadConfig(sceneName string, params map[string]interface{}) *ParamsMapList {
	pl.Params = make(map[string]string)
	pl.SceneName = sceneName
	// 初始化
	pl.NestedGroup = make(map[string]map[string]bool)
	pl.OriginList = make(map[string]*SearchCondition)
	pl.generateConfig(params)
	pl.GenerateEsQueryConfig()
	return pl
}

// 根据传参, 生成原始 config 信息
func (pl *ParamsMapList) generateConfig(params map[string]interface{}) *ParamsMapList {
	for k, v := range params {
		// 这里 这个 key 已经是es 里配置好的字段了, 包括带 .  的这种二级字段
		key, exist := pl.Params[k]
		if !exist {
			// 没有配置的 key 就按默认的一级走
			key = k
		}
		raw := pl.generateOrigin(key, v)
		if raw == nil {
			logs.Trace("返回空值  key: %v, v: %v", key, v)
			continue
		}
		pl.OriginList[key] = raw
	}

	return pl
}

func (pl *ParamsMapList) generateOrigin(key string, value interface{}) *SearchCondition {
	karr := strings.Split(key, ".")

	tmp := &SearchCondition{ // 默认是用 must
		Type:  QueryMust,
		Key:   key,
		Value: value,
	}
	// nested 类型, 需要备注 path 地址
	var path string
	if len(karr) > 1 {
		path = karr[0]
	} else {
		path = TypeTerm
	}

	if _, ok := pl.NestedGroup[path]; !ok {
		pl.NestedGroup[path] = map[string]bool{
			key: true,
		}
	} else {
		pl.NestedGroup[path][key] = true
	}
	return tmp
}

func (pl *ParamsMapList) GenerateEsQueryConfig() []*EsQueryCondition {
	result := make([]*EsQueryCondition, 0)
	if pl.NestedGroup == nil || len(pl.NestedGroup) == 0 || len(pl.OriginList) == 0 {
		return result
	}

	// 初始化存在
	for pth := range pl.NestedGroup {
		var tmp *EsQueryCondition
		// 一级 term 字段
		if pth == TypeTerm {
			tmp = &EsQueryCondition{
				Type:       TypeTerm,
				Path:       "",
				Conditions: make([]*SearchCondition, 0),
			}
			tmp.LoadConditionsFromParamsConfig(pl.NestedGroup[pth], pl)
		} else {
			// 子文档字段
			tmp = &EsQueryCondition{
				Type:       TypeNested,
				Path:       pth,
				Conditions: make([]*SearchCondition, 0),
			}
			tmp.LoadConditionsFromParamsConfig(pl.NestedGroup[pth], pl)
		}

		result = append(result, tmp)
	}

	pl.EsQueryConfig = &EsQueryConfig{
		Conditions: result,
	}

	return pl.EsQueryConfig.Conditions
}

func (pl *ParamsMapList) GenerateQuery() elastic.Query {
	return pl.EsQueryConfig.BuildBoolQuery()
}

// GetKeyBodyOnce 从 OriginList 去除
func (pl *ParamsMapList) GetKeyBodyOnce(path string, key string) (item *SearchCondition, exist bool) {
	item, exist = pl.OriginList[key]
	if !exist {
		return
	}
	delete(pl.OriginList, key)

	tmp := pl.NestedGroup[path]
	delete(tmp, key)
	// 移除已经用上的配置
	pl.NestedGroup[path] = tmp
	if len(tmp) == 0 {
		delete(pl.NestedGroup, path)
	}
	return
}

// 判断是否是多个参数
func getMultiTerms(value interface{}) (terms []interface{}, isMulti bool) {
	ttValue, isArr := value.([]interface{})
	if isArr {
		return ttValue, true
	}

	terms = make([]interface{}, 0)
	tmpValue, isString := value.(string)
	if !isString {
		terms = append(terms, value)
		return terms, false
	}

	valueArr := strings.Split(tmpValue, ",")
	for cur := range valueArr {
		terms = append(terms, valueArr[cur])
	}

	if len(valueArr) > 1 {
		isMulti = true
	}
	return terms, isMulti
}

// 检测是否是范围类型的数值
func checkRangeValue(key string, value interface{}) (rangeQuery *elastic.RangeQuery, isRange bool) {
	valueArr, ok := value.([]interface{})
	if ok && len(valueArr) == 2 { // 如果恰好是两个值得数组, 则说明是range 类型的
		isRange = true
	} else {
		return nil, false
	}
	query := elastic.NewRangeQuery(key)
	if valueArr[0] != nil {
		query.Gte(valueArr[0])
	}

	if valueArr[1] != nil {
		query.Lte(valueArr[1])
	}
	return query, isRange
}

func checkNotNullValue(key string, value interface{}) bool {
	vStr, ok := value.(string)
	if !ok {
		return false
	}
	var isNotNull bool
	if vStr == SignNotNull {
		isNotNull = true
	}
	return isNotNull
}

func checkSubKey(key string) string {
	strArr := strings.Split(key, ">")

	if len(strArr) > 0 {
		return strings.Join(strArr, ".")
	}
	return key
}

// 不作为 query 参数使用, 而是放在了agg指标当中
func ignoreColumns(key string) bool {
	if key == SignTimeRange {
		return true
	}
	return false
}

func (ec *EsQueryCondition) BuildWithQuery(query *elastic.BoolQuery) *elastic.BoolQuery {
	if len(ec.Conditions) == 0 {
		return query
	}

	tmpQuery := elastic.NewBoolQuery()
	if ec.Type != TypeNested {
		tmpQuery = query
	}

	maxCnt := len(ec.Conditions)
	jumpCnt := 0

	for i := 0; i < len(ec.Conditions); i++ {
		cur := ec.Conditions[i]

		// 过滤筛选字段
		if ignoreColumns(cur.Key) {
			jumpCnt += 1
			continue
		}

		if cur.Optional > 0 && cur.Value == nil {
			//logs.Trace("空的可选参数赋值 key: %v", cur.Key)
			jumpCnt += 1
			continue
		}

		terms, isMulti := getMultiTerms(cur.Value)

		rangeQuery, isRange := checkRangeValue(cur.Key, cur.Value)

		logs.Trace("cur.Key: %v, isRange:%v", cur.Key, isRange)

		if isRange {
			cur.Type = QueryRange
		}

		// 检测是否要求非空
		isNotNull := checkNotNullValue(cur.Key, cur.Value)
		if isNotNull {
			cur.Type = QueryMustNot
			cur.Value = ""
		}

		// 判断是否是二级字段 (这就是关键字符冲突造的孽)
		cur.Key = checkSubKey(cur.Key)

		switch cur.Type {
		case QueryRange:
			tmpQuery = tmpQuery.Filter(rangeQuery)
		case QueryMust:
			if isMulti {
				tmpQuery = tmpQuery.Filter(elastic.NewTermsQuery(cur.Key, terms...))
			} else {
				tmpQuery = tmpQuery.Filter(elastic.NewTermQuery(cur.Key, cur.Value))
			}
		case QueryMustMulti:
			v := strings.Split(cur.Value.(string), ",")
			tmpI := make([]interface{}, 0)
			for _, vv := range v {
				tmpI = append(tmpI, vv)
			}
			tmpQuery = tmpQuery.Filter(elastic.NewTermsQuery(cur.Key, tmpI...))
		case QueryMustNot:
			if isMulti {
				tmpQuery = tmpQuery.MustNot(elastic.NewTermsQuery(cur.Key, terms...))
			} else {
				tmpQuery = tmpQuery.MustNot(elastic.NewTermQuery(cur.Key, cur.Value))
			}
		case QueryMustNotMulti:
			v := strings.Split(cur.Value.(string), ",")
			tmpI := make([]interface{}, 0)
			for _, vv := range v {
				tmpI = append(tmpI, vv)
			}
			tmpQuery = tmpQuery.MustNot(elastic.NewTermsQuery(cur.Key, tmpI...))
		default:
			// 默认走 term
			tmpQuery = tmpQuery.Filter(elastic.NewTermQuery(cur.Key, cur.Value))
		}
	}

	if ec.Type == TypeNested {
		// 说明没有完全略过, 里面的参数还要继续用
		if jumpCnt < maxCnt {
			q := elastic.NewNestedQuery(ec.Path, tmpQuery)
			query.Filter(q)
		}
	}

	return query
}

func (ecg *EsQueryConfig) BuildBoolQuery() *elastic.BoolQuery {
	query := elastic.NewBoolQuery()

	if ecg == nil {
		return query
	}

	if len(ecg.Conditions) == 0 {
		return query
	}

	for i := 0; i < len(ecg.Conditions); i++ {
		query = ecg.Conditions[i].BuildWithQuery(query)
	}

	return query
}

// CombineParamsConfig 将 params 的配置合并到当前搜索中
func (eqc *EsQueryConfig) CombineParamsConfig(paramsConfig *ParamsMapList) {
	for i := 0; i < len(eqc.Conditions); i++ {
		current := eqc.Conditions[i]
		path := current.Path

		// 获取空值
		if path == "" {
			path = TypeTerm
		}
		keyGroup, ok := paramsConfig.NestedGroup[path]
		if !ok {
			continue
		}

		// 如果配置中已经存在值, 就给去掉
		for j := 0; j < len(current.Conditions); j++ {
			tmpItem := current.Conditions[j]
			paramsConfig.GetKeyBodyOnce(path, tmpItem.Key)
		}

		for key := range keyGroup {
			keyItem, exist := paramsConfig.GetKeyBodyOnce(path, key)
			if !exist {
				continue
			}
			eqc.Conditions[i].Conditions = append(eqc.Conditions[i].Conditions, keyItem)
		}
	}

	// 把最后没配置上但是传进来的参数给拼到后面去
	tmpList := paramsConfig.GenerateEsQueryConfig()
	eqc.Conditions = append(eqc.Conditions, tmpList...)
	logs.Trace("CombineParamsConfig: %v", eqc.Conditions)
}

// LoadValueFromParams 从外部传入的 params 中获取参数
func (eqc *EsQueryConfig) LoadValueFromParams(params map[string]interface{}) error {
	//logs.Trace("LoadValueFromParams")
	if eqc.Conditions == nil || len(eqc.Conditions) == 0 {
		return nil
	}
	//logs.Trace("LoadValueFromParams cond: %v", eqc.Conditions)

	for i := 0; i < len(eqc.Conditions); i++ {
		err := eqc.Conditions[i].LoadValueFromParams(params)
		if err != nil {
			return err
		}
	}

	return nil
}

// LoadValueFromParams 从外部传入的 params 中获取参数
func (ec *EsQueryCondition) LoadValueFromParams(params map[string]interface{}) error {
	if ec.Conditions == nil || len(ec.Conditions) == 0 {
		return nil
	}
	for i := 0; i < len(ec.Conditions); i++ {
		err := ec.Conditions[i].LoadValueFromParams(params)
		// 可选参数过滤
		if err != nil && ec.Conditions[i].Optional > 0 {
			continue
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// LoadConditionsFromParamsConfig 从参数配置中加载对应的配置
func (ec *EsQueryCondition) LoadConditionsFromParamsConfig(mapList map[string]bool, pl *ParamsMapList) {
	if len(mapList) == 0 {
		return
	}
	tmp := make([]*SearchCondition, 0)

	for key := range mapList {
		item, ok := pl.OriginList[key]
		if !ok {
			continue
		}
		tmp = append(tmp, item)
	}
	ec.Conditions = tmp
}

// LoadValueFromParams 从外部传入的 params 中获取参数
func (sc *SearchCondition) LoadValueFromParams(params map[string]interface{}) error {
	karr := strings.Split(sc.Key, ".")
	//logs.Trace("LoadValueFromParams key: %v, karr: %v", sc.Key, karr)
	if len(karr) == 0 || len(karr) > 2 {
		return fmt.Errorf("key 命名错误: %v", sc.Key)
	}
	if sc.Value != nil {
		logs.Trace("已存在 value, 不再重复赋值, key: %v", sc.Key)
		return nil
	}
	key := karr[len(karr)-1]
	pv, ok := params[key]
	if !ok {
		// 如果不存在, 就变成可选参数, 不再报错
		sc.Optional = 1
		return nil
		//return fmt.Errorf("缺少执行搜索的必要参数, 参数key: %v", key)
	}
	logs.Trace("赋值 key: %v, value: %v", sc.Key, pv)
	// 将
	sc.Value = pv
	return nil
}
