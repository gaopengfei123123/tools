package mock

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

var cityPool *RandPool            // 城市
var itemPool map[string]*RandPool // 类目
var genderPool *RandPool          // 性别
var orderPool *RandPool           // 近30天是否下单
var osPool *RandPool              // 手机平台
var sourcePagePool *RandPool      // 来源页
var remotePool *RandPool          // 是否异地
var hourPool *RandPool            // 时间段

type RandPool struct {
	OriList map[int]string    // 获取原始的配比
	TmpList map[string][2]int // 获取生成的 map
	Total   int
}

func (rp *RandPool) LoadConfig(data map[string]int) *RandPool {
	if data == nil {
		return rp
	}

	rp.TmpList = make(map[string][2]int)
	for v, k := range data {
		curStart := rp.Total
		rp.Total += k
		curEnd := rp.Total
		tmp := [2]int{curStart, curEnd}
		rp.TmpList[v] = tmp
	}
	return rp
}
func (rp *RandPool) GetItem() string {
	curV := GetRandInt(rp.Total)

	for item, rng := range rp.TmpList {
		start := rng[0]
		end := rng[1]
		if curV >= start && curV < end {
			return item
		}
	}
	return "未知"
}

func GetRandInt(max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max)
}

func GetRandCity(exclude ...string) string {
	if cityPool == nil {
		mp := map[string]int{
			"北京": 100,
			"上海": 80,
			"杭州": 50,
			"广东": 50,
			"深圳": 45,
			"重庆": 40,
			"成都": 30,
		}

		cityPool = new(RandPool)
		cityPool.LoadConfig(mp)
	}

	if len(exclude) != 0 {
		v := exclude[0]
		for {
			ct := cityPool.GetItem()
			if ct != v {
				return ct
			}
		}
	}

	return cityPool.GetItem()
}

// GetRandItem 获取随机类目
func GetRandItem(city string) string {
	if itemPool == nil {
		itemPool = make(map[string]*RandPool)

		conf := map[string]int{
			"紧致抗衰": 1454,
			"除皱瘦脸": 1173,
			"玻尿酸":  953,
			"吸脂":   606,
			"美白嫩肤": 586,
			"抗敏修复": 411,
			"鼻综合":  432,
			"眼综合":  405,
			"面部提升": 400,
			"隆胸":   305,
		}
		rp := new(RandPool)
		itemPool["通用"] = rp.LoadConfig(conf)
	}

	if cp, exist := itemPool[city]; exist {
		return cp.GetItem()
	} else {
		return itemPool["通用"].GetItem()
	}
}

// GetRandGender 获取随机性别
func GetRandGender() string {
	if genderPool == nil {
		conf := map[string]int{
			"男性":  3,
			"女性":  7,
			"未知":  1,
			"未说明": 3,
		}
		genderPool = new(RandPool)
		genderPool.LoadConfig(conf)
	}
	return genderPool.GetItem()
}

// GetRandOrder 获取随机订单类型
func GetRandOrder() string {
	if orderPool == nil {
		conf := map[string]int{
			"下单":  300,
			"未下单": 1000,
		}
		orderPool = new(RandPool)
		orderPool.LoadConfig(conf)
	}
	return orderPool.GetItem()
}

// GetRandOS 获取随机平台
func GetRandOS() string {
	if osPool == nil {
		conf := map[string]int{
			"IOS":     61,
			"Android": 39,
		}
		osPool = new(RandPool)
		osPool.LoadConfig(conf)
	}
	return osPool.GetItem()
}

// GetRandRemote  是否异地
func GetRandRemote() string {
	if remotePool == nil {
		conf := map[string]int{
			"本地": 6,
			"异地": 4,
		}
		remotePool = new(RandPool)
		remotePool.LoadConfig(conf)
	}
	return remotePool.GetItem()
}

// GetRandHour  获取随机的时间段
func GetRandHour() int {
	if hourPool == nil {
		conf := map[string]int{
			"0":  3,
			"1":  2,
			"2":  1,
			"3":  1,
			"4":  1,
			"5":  1,
			"6":  1,
			"7":  1,
			"8":  2,
			"9":  3,
			"10": 4,
			"11": 5,
			"12": 5,
			"13": 6,
			"14": 5,
			"15": 4,
			"16": 3,
			"17": 3,
			"18": 4,
			"19": 5,
			"20": 6,
			"21": 7,
			"22": 6,
			"23": 4,
		}
		hourPool = new(RandPool)
		hourPool.LoadConfig(conf)
	}
	hh := hourPool.GetItem()
	v, _ := strconv.Atoi(hh)
	return v
}

// GetRandSourcePage 获取来源页   首页开屏,首页feed,搜索结果页,类目页,其他
func GetRandSourcePage() string {
	if sourcePagePool == nil {
		conf := map[string]int{
			"首页首屏":   40,
			"首页feed": 40,
			"搜索结果页":  100,
			"类目页":    60,
			"其他":     40,
		}
		sourcePagePool = new(RandPool)
		sourcePagePool.LoadConfig(conf)
	}
	return sourcePagePool.GetItem()
}

func GetRandIpaddr4() string {
	ip := fmt.Sprintf("%d.%d.%d.%d", GetRandInt(255), GetRandInt(255), GetRandInt(255), GetRandInt(255))
	return ip
}

func GetRandDeviceID() int {
	return 1000000 + GetRandInt(9000000)
}
