package sort

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/gaopengfei123123/tools/convert"
	"sort"
	"strings"
)

// SortedList 排序用数组
type SortedList struct {
	List       []map[string]interface{} // 数据本体
	column     string                   // 要排序的字段
	ColumnType string                   // 要排序的类型
	orderDesc  bool                     // 顺序  true:升序,false:降序
}

// LoadData 加载数据
func (sor *SortedList) LoadData(data interface{}) *SortedList {
	bt, _ := json.Marshal(data)
	//logs.Info("sortList: %s err: %v", bt, err)

	tmp := make([]map[string]interface{}, 0)
	_ = json.Unmarshal(bt, &tmp)
	sor.List = tmp
	return sor
}

func (sor *SortedList) GroupSorted(group string, column string, orderDesc bool) *SortedList {
	// 存储分组的顺序
	groupList := make(map[string]int)

	// 实际的分组
	tmpListGroup := make([][]map[string]interface{}, len(sor.List))

	//debug.PrintJson("sortList", sor.List)

	// 遍历整个 list, 根据 group 字段来拆分数组
	for i := 0; i < len(sor.List); i++ {
		current := sor.List[i]

		groupValue, ok := current[group]
		groupValueStr := fmt.Sprintf("%v", groupValue)
		if !ok {
			logs.Error("groupValueStrErr")
			continue
		}

		_, ok = groupList[groupValueStr]
		if !ok {
			groupIndex := len(groupList)
			groupList[groupValueStr] = groupIndex
		}

		curIndex := groupList[groupValueStr]
		tmpListGroup[curIndex] = append(tmpListGroup[curIndex], current)
	}

	resultList := make([]map[string]interface{}, 0, len(sor.List))

	for i := 0; i < len(tmpListGroup); i++ {
		curList := tmpListGroup[i]

		tmpSortList := new(SortedList)
		tmpSortList.LoadData(curList).Sorted(column, orderDesc)

		//debug.PrintJson("tmpListSort", tmpSortList)

		//for j := 0; j < len(tmpSortList.List); j++ {
		//	resultList = append(resultList, tmpSortList.List[j])
		//}
		resultList = append(resultList, tmpSortList.List...)
	}

	b, _ := convert.JSONEncode(tmpListGroup)
	logs.Debug("groupList: %s", b)

	sor.List = resultList
	return sor
}

// Sorted 根据指定的字段和排序规则进行排序
func (sor *SortedList) Sorted(column string, orderDesc bool) *SortedList {
	sor.column = column
	sor.orderDesc = orderDesc
	sort.Sort(sor)
	return sor
}

// DecodeList 将列表赋值出去, 需要传递指针类型参数
func (sor *SortedList) DecodeList(data interface{}) error {
	bt, _ := json.Marshal(sor.List)
	err := json.Unmarshal(bt, data)
	return err
}

// Len 排序需要实现的接口
func (sor *SortedList) Len() int {
	return len(sor.List)
}

// Swap 排序需要实现的接口
func (sor *SortedList) Swap(i, j int) {
	sor.List[i], sor.List[j] = sor.List[j], sor.List[i]
}

// Less 排序需要实现的接口
func (sor *SortedList) Less(i, j int) bool {
	if sor.column == "" {
		return false
	}

	if sor.Len() == 0 {
		return false
	}

	// 默认走 float64, 有特殊情况再说
	k1 := sor.getValueByKey(sor.List[i], sor.column)
	k2 := sor.getValueByKey(sor.List[j], sor.column)

	if sor.orderDesc {
		return k1 > k2
	} else {
		return k1 < k2
	}
}

func (sor *SortedList) getValueByKey(raw map[string]interface{}, key string) float64 {
	l := strings.Split(key, ".")
	// 一维对象, 直接取对应 key 的值
	if len(l) == 1 {
		v, ok := raw[key]
		if !ok {
			return 0
		}
		value, ok := v.(float64)
		if !ok {
			return 0
		}
		return value
	}

	if len(l) != 2 {
		return 0
	}

	// 下面只处理二维的, 不作递归了

	key1, key2 := l[0], l[1]
	raw1, ok := raw[key1]
	if !ok {
		return 0
	}

	raw2, ok := raw1.(map[string]interface{})
	if !ok {
		return 0
	}

	v2, ok := raw2[key2]
	if !ok {
		return 0
	}

	value2, ok := v2.(float64)
	if !ok {
		return 0
	}

	return value2
}
