package tools

import (
	"context"
	"fmt"
	"github.com/gaopengfei123123/tools/timeutil"
	"github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
	"time"
)

// GetESMapping 获取指定索引的 mapping
func GetESMapping(ctx context.Context, client *elastic.Client, index ...string) (result map[string]interface{}, err error) {
	if ctx == nil {
		ctx = context.TODO()
	}
	tmpIndexInfo, err := client.GetMapping().Index(index...).Do(ctx)
	result = tmpIndexInfo
	return
}

func CheckEsMapExist(ctx context.Context, client *elastic.Client, index string) bool {
	if ctx == nil {
		ctx = context.TODO()
	}
	res, _ := client.CatIndices().Index(index).Do(ctx)
	// 已经创建过索引就不再创建
	if len(res) > 0 && res[0].Status == "open" {
		return true
	}
	return false
}

// CopyEsIndex 复制 es 索引
func CopyEsIndex(ctx context.Context, client *elastic.Client, src string, target string) (success bool, err error) {
	if ctx == nil {
		ctx = context.TODO()
	}

	srcExist := CheckEsMapExist(ctx, client, src)
	targetExist := CheckEsMapExist(ctx, client, target)

	// 源数据不存在
	if !srcExist {
		success = false
		err = errors.Errorf("%v", "源索引不存在")
		return
	}

	// 已经存在索引
	if targetExist {
		success = true
		err = errors.Errorf("%v", "目标索引已存在, 不再处理")
		return
	}

	mp, err := GetESMapping(ctx, client, src)
	if err != nil {
		success = false
		err = errors.Errorf("Get Mapping Err:%v", err)
	}

	var mapping interface{}
	if info, ok := mp[src]; ok {
		mapping = info
	}

	_, err = client.CreateIndex(target).BodyJson(mapping).Do(ctx)
	if err != nil {
		err = errors.Errorf("Create Index Err:%v", err)
	}
	success = true
	return
}

// GetEarliestDate 分片查询时间临界点
func GetEarliestDate() time.Time {
	// TODO 更改线上索引分界时间, 测试环境是  2023-05-15
	earliestDate, _ := timeutil.StrToTime("2023-05-31 22:10:00", time.Local)
	return earliestDate
}

// GetCurrentSliceESIndex 分片基础函数
func GetCurrentSliceESIndex(indexName string, timer ...time.Time) (index string) {
	var tt time.Time
	if len(timer) == 0 {
		tt = time.Now()
	} else {
		tt = timer[0]
	}

	// TODO  低于这个日期的, 返回老索引
	if tt.Unix() < GetEarliestDate().Unix() {
		return indexName
	}
	yyyymm := tt.Format("200601")
	return fmt.Sprintf("%s_%s", indexName, yyyymm)
}

// GetSliceESIndexByTimeRange 根据时间格式过滤索引
func GetSliceESIndexByTimeRange(indexName string, startTime time.Time, endTime time.Time) (indexArr []string, err error) {
	indexArr = make([]string, 0, 10)
	// 重复索引过滤
	existMap := make(map[string]struct{})

	end := endTime.AddDate(0, 1, -1)
	for startTime.Unix() < end.Unix() {
		//logs.Info("\n")
		//logs.Info("start:%v", startTime)
		//logs.Info("earliestTime: %v", GetEarliestDate().Format("200601"))
		//logs.Info("curTime: %v", startTime.Format("200601"))
		// 如果是和分界时间同一个月的, 就再加个当月的索引
		if GetEarliestDate().Format("200601") == startTime.Format("200601") {
			//logs.Info("special start:%v", startTime)
			yyyymm := startTime.Format("200601")
			tmpp := fmt.Sprintf("%s_%s", indexName, yyyymm)
			if _, exist := existMap[tmpp]; !exist {
				indexArr = append(indexArr, tmpp)
				existMap[tmpp] = struct{}{}
			}

			// 当月出现两个索引
			tmpp = indexName
			if _, exist := existMap[tmpp]; !exist {
				indexArr = append(indexArr, tmpp)
				existMap[tmpp] = struct{}{}
			}
		}

		tmp := GetCurrentSliceESIndex(indexName, startTime)
		if _, ok := existMap[tmp]; ok {
			//logs.Info("jump %v", startTime)
			startTime = startTime.AddDate(0, 0, 27)
			//logs.Info("jump TO %v", startTime)
			continue
		}
		existMap[tmp] = struct{}{}

		// 最多暴露超过当前一个月的索引
		if startTime.Year() >= time.Now().Year() && startTime.Month() > time.Now().Month()+1 {
			break
		}

		indexArr = append(indexArr, tmp)
		startTime = startTime.AddDate(0, 0, 27)
	}

	// 如果传递的时间有错误, 兜底使用startDate索引
	if len(indexArr) == 0 {
		tmp := GetCurrentSliceESIndex(indexName, startTime)
		indexArr = append(indexArr, tmp)
	}
	return
}

// GetSliceESIndexByDateRange 根据时间格式过滤索引 可传递的格式有 2006-01-02 2006-1-02 2006-1-2 20060102
func GetSliceESIndexByDateRange(indexName string, startDate string, endDate string) (indexArr []string, err error) {
	tmpResult := []string{GetCurrentSliceESIndex(indexName)}
	startTime, e := timeutil.StrToTime(startDate, time.UTC)
	if e != nil {
		return tmpResult, e
	}

	endTime, e := timeutil.StrToTime(endDate, time.UTC)
	if e != nil {
		return tmpResult, e
	}
	return GetSliceESIndexByTimeRange(indexName, startTime, endTime)
}

// GetEsIndex 默认获取批量索引的方法
func GetEsIndex(indexName string, params map[string]interface{}) []string {
	result := make([]string, 0, 1)

	// 默认返回前三个月的索引, 如果当前是27号后, 会带上未来一个月的索引
	var startDate, endDate string
	if i, exist := params["start_date"]; exist {
		startDate = fmt.Sprintf("%v", i)
	}
	if i, exist := params["end_date"]; exist {
		endDate = fmt.Sprintf("%v", i)
	}

	if startDate == "" && endDate == "" {
		start := time.Now().AddDate(0, -3, 0)
		end := time.Now()

		if end.Day() >= 27 {
			end = end.AddDate(0, 1, 0)
		}
		indexList, err := GetSliceESIndexByTimeRange(indexName, start, end)
		if err != nil {
			index := GetCurrentSliceESIndex(indexName)
			result = append(result, index)
			return result
		}

		result = append(result, indexList...)
		return result
	}

	// 如果存在明确的时间范围, 则按时间范围的来
	result, _ = GetSliceESIndexByDateRange(indexName, startDate, endDate)
	return result
}
