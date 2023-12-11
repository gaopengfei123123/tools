package trie

import (
	"strings"
	"sync"
)

type Null struct{}

// InvalidWords 不参与敏感词判断
const InvalidWords = " ,~,!,@,#,$,%,^,&,*,(,),_,-,+,=,?,<,>,—,，,。,/,\\,|,《,》,？,;,:,：,',‘,；,“,"

// GetInvalidWordMap 获取非敏感词列表
func GetInvalidWordMap() map[string]Null {
	words := strings.Split(InvalidWords, ",")
	result := make(map[string]Null)
	for _, v := range words {
		result[v] = Null{}
	}
	return result
}

type Trie struct {
	BlackList []string `gorm:"-" json:"black_list"` // 从黑名单表里读到的敏感词集合
	TrieTree  *Twig    `json:"trie_tree"`           //单词查找树
	sync.Mutex
}

// BlackListToTrieTree 把黑名单里的词转化成树状结构
func (rule *Trie) BlackListToTrieTree() *Trie {
	if len(rule.BlackList) == 0 {
		return rule
	}

	// 加锁
	rule.Lock()
	defer rule.Unlock()

	set := make(map[string]Null)
	for _, v := range rule.BlackList {
		set[v] = Null{}
	}

	rule.AddSensitiveToTree(set)

	rule.BlackList = make([]string, 0, 1)
	return rule
}

// AddSensitiveToTree 批量将敏感词加入字典中
func (rule *Trie) AddSensitiveToTree(set map[string]Null) {
	if rule.TrieTree == nil {
		rule.TrieTree = NewTrieTwig()
	}

	for key := range set {
		str := []rune(key)
		nowMap := rule.TrieTree
		ln := len(str)
		for i := 0; i < ln; i++ {
			curStr := string(str[i])
			if _, ok := nowMap.ChildMap[curStr]; !ok { //如果该key不存在，
				thisMap := NewTrieTwig()
				thisMap.RootDeep = ln - i - 1
				thisMap.HeadHeight = i + 1
				nowMap.ChildMap[curStr] = thisMap
				nowMap = thisMap
			} else {
				nowMap = nowMap.ChildMap[curStr]
				// 如果有更长的分支, 则更新当前根的统计
				curLen := ln - i - 1
				if curLen > nowMap.RootDeep {
					nowMap.RootDeep = curLen
				}
			}

			// 标记出来是一个结尾
			if i == len(str)-1 {
				nowMap.IsEnd = true
			}
		}
	}
}

// SignSensitiveWords 在原文字中标记出来命中的敏感词
func (rule *Trie) SignSensitiveWords(rawText string, startTag, endTag string) (result string, hitKey []string) {
	// 加锁
	rule.Lock()
	defer rule.Unlock()

	str := []rune(rawText)
	// 获取需要忽略检测的字符(白名单)
	invalidWord := GetInvalidWordMap()

	hitKeyWord := make(map[string]Null)
	// 命中节点的起止位置,结构 =>   起始:结束
	hitKeyWordCoord := make(map[int]int)

	for i := 0; i < len(str); i++ {
		cur := str[i]
		curStr := string(cur)
		//logs.Info("i: %d, cur: %v, curStr: %v", i, cur, curStr)
		if _, ok := invalidWord[curStr]; ok || curStr == "," {
			continue
		}

		// 读到一个字和关键词中匹配的, 开始校验是否存在命中词, 直到和这个字关联的所有的词都扫出来, 再移到下一个字
		if thisTwig, ok := rule.TrieTree.ChildMap[curStr]; ok {
			//debug.PrintJson("thisTwig", thisTwig)
			tmpLen := i + thisTwig.RootDeep + 1
			if tmpLen > len(str) {
				tmpLen = len(str)
			}

			keyWordList := thisTwig.GetHitKeywords(str[i:tmpLen])
			//debug.PrintJson("keywordList", keyWordList, true)
			if len(keyWordList) != 0 {
				maxLen := 0
				for _, k := range keyWordList {
					hitKeyWord[k] = Null{}
					keyL := len([]rune(k))
					if maxLen < keyL {
						maxLen = keyL
					}
				}
				hitKeyWordCoord[i] = maxLen + i
			}
		}
	}

	//debug.PrintJson("hitKeyword", hitKeyWord, true)
	//debug.PrintJson("hitKeyWordCoord", hitKeyWordCoord, true)
	//logs.Info("curStr: %s", string(str[3:10]))

	hitKey = make([]string, 0, len(hitKeyWord))

	// 不包含高亮标签时的导出
	if startTag == "" && endTag == "" {
		for cur := range hitKeyWord {
			hitKey = append(hitKey, cur)
		}
		return rawText, hitKey
	}

	// 替换逻辑
	if len(hitKeyWord) > 0 && startTag != "" && endTag != "" {
		start := []rune(startTag)
		end := []rune(endTag)
		// 这里把高亮标签占用的长度也算上, 减少内存分配次数
		resultArr := make([]rune, 0, len(str)+len(hitKeyWord)*(len(start)+len(end)))
		for i := 0; i < len(str); i++ {
			// 当前下标存在关键词, 需要替换
			if endIndex, exist := hitKeyWordCoord[i]; exist {
				// 加上高亮前缀
				resultArr = append(resultArr, start...)
				// 加文字本体
				resultArr = append(resultArr, str[i:endIndex]...)
				// 加高亮后缀
				resultArr = append(resultArr, end...)
				// 跳转到文字后将遍历指针放到替换的文字后, 继续查询新的数据
				i = endIndex - 1
				//logs.Info("resultArr i: %v Cap: %v", i, cap(resultArr)) // 校验内存不再分配
				continue
			}
			resultArr = append(resultArr, str[i])
			//logs.Info("resultArr i: %v Cap: %v", i, cap(resultArr))
		}
		result = string(resultArr)

		// 汇总去重后的关键词
		for key := range hitKeyWord {
			hitKey = append(hitKey, key)
		}

		return
	}

	return rawText, hitKey
}

type Twig struct {
	IsEnd      bool             `json:"is_end"`
	RootDeep   int              `json:"root_deep"`   // 当前分支最深的层级
	HeadHeight int              `json:"head_height"` // 距离头部最高的层级
	ChildMap   map[string]*Twig `json:"child_map"`
}

func NewTrieTwig() *Twig {
	tmp := new(Twig)
	tmp.IsEnd = false
	tmp.ChildMap = make(map[string]*Twig)
	return tmp
}

func (twig *Twig) GetHitKeywords(str []rune) []string {
	result := make([]string, 0, 1)

	// 当前节点命中关键字
	if twig.IsEnd {
		result = append(result, string(str[:twig.HeadHeight]))
	}

	// 到达最底端, 不再处理下一层 / 已经到节点末端, 不再执行后续
	if len(str) == twig.HeadHeight || twig.RootDeep == 0 {
		return result
	}

	// 之后循环时对应的字
	nextStr := string(str[twig.HeadHeight])

	// 查找后续节点
	childMap, exist := twig.ChildMap[nextStr]
	if !exist {
		return result
	}

	//logs.Info("curStr: %s", string(str))
	//logs.Info("curOffset: %v", twig.HeadHeight)
	//logs.Info("curWord: %s", first)
	//debug.PrintJson("curHead", childMap)

	// 获取后续节点数据
	keywords := childMap.GetHitKeywords(str)
	result = append(result, keywords...)
	return result
}
