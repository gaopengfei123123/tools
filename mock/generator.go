package mock

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/gocarina/gocsv"
	"math"
	"os"
	"time"
)

func MockDataGenerator() {
	logs.Trace("xxx")
	return
}

func OutputCsvDemo() {
	f, err := os.Create("demo.csv")
	//关闭流
	defer f.Close()
	//写入UTF-8 格式
	f.WriteString("\xEF\xBB\xBF")
	//var newContent []CSVItem
	list := GenerateNormalGuy()
	//newContent = make([]CSVItem, 0, 1)
	//添加数据
	//newContent = append(newContent, CSVItem{ID: 1, IP: "GPF"})
	//newContent = append(newContent, CSVItem{ID: 2, IP: "GPF"})
	//保存文件流
	err = gocsv.MarshalFile(list, f)
	if err != nil {
		logs.Trace("err: %v", err)
		logs.Trace("end")
		return
	}
}

func OutPutCsv(total int) {
	f, err := os.OpenFile("pv_data_fake.csv", os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	//关闭流
	defer f.Close()
	//写入UTF-8 格式
	f.WriteString("\xEF\xBB\xBF")

	res := make([]CSVItem, 0, total*10)
	for i := 0; i < total; i++ {
		//list := GenerateNormalGuy()
		list := GenerateFakeGuy()

		res = append(res, list...)
	}

	err = gocsv.Marshal(res, f)

	if err != nil {
		logs.Trace("err: %v", err)
		logs.Trace("end")
		return
	}
}

func OutPutCsvSummary(total int) {
	f, err := os.OpenFile("pv_data_fake_summary.csv", os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	//关闭流
	defer f.Close()
	//写入UTF-8 格式
	f.WriteString("\xEF\xBB\xBF")

	res := make([]CSVItemSummary, 0, total)
	for i := 0; i < total; i++ {
		//list := GenerateNormalGuy()
		item := GenerateFakeSummary()

		res = append(res, item)
	}

	err = gocsv.Marshal(res, f)

	if err != nil {
		logs.Trace("err: %v", err)
		logs.Trace("end")
		return
	}
}

type CSVItem struct {
	ID          int
	IP          string
	DeviceID    string
	City        string // 用户所在地
	Item        string // 商品类目
	Gender      string // 性别
	Order       string // 是否产生订单
	OS          string // 手机系统
	SourcePage  string // 点击来源页
	ProductCity string // 商品所在地
	IsRemote    string // 异地/本地
	CreateTime  string // 点击时间
	Hour        string // 所属小时
	Interval    string // 距离上次点击的间隔
}

// CSVItemSummary 单条数据汇总
type CSVItemSummary struct {
	ID          int
	IP          string
	DeviceID    string
	City        string // 用户所在地
	Item        string // 商品类目
	Gender      string // 性别
	Order       string // 是否产生订单
	OS          string // 手机系统
	SourcePage  string // 点击来源页
	ProductCity string // 商品所在地
	IsRemote    string // 异地/本地
	CreateTime  string // 点击时间
	Hour        string // 所属小时
	Interval    string // 距离上次点击的间隔
	TotalCnt    int    // 总点击数
	ItemCnt     int    // 点击类目数
	AVGTime     int    // 平均点击间隔
	TotalTime   int    // 总点击耗时
	IPCnt       int    // 相同 id的 ip 数
	AvgItemCnt  int    // 平均点击类目数
	IsRisk      int    // 风控标识
	RiskType    int    // 风控类型
}

// GenerateNormalGuy 生成正常的日志
func GenerateNormalGuy() []CSVItem {
	// 总条数
	num := GetRandInt(20)

	userID := GetRandInt(99999999)
	IP := GetRandIpaddr4()
	DeviceID := GetRandDeviceID()
	city := GetRandCity(false)
	productCity := city
	isRemote := GetRandRemote()
	order := GetRandOrder()
	iphoneOS := GetRandOS(false)
	sourcePage := GetRandSourcePage()
	gender := GetRandGender()

	if isRemote == "外地" {
		productCity = GetRandCity(false, city)
	}

	StartTime := time.Date(2023, 02, 1+GetRandInt(3), GetRandHour(false), GetRandInt(59), 0, 0, time.Local)

	result := make([]CSVItem, 0, num)
	for i := 0; i < num; i++ {
		randSec := 2 + GetRandInt(30)
		StartTime = StartTime.Add(time.Second * time.Duration(randSec))

		tmp := CSVItem{
			ID:          userID,
			IP:          IP,
			DeviceID:    fmt.Sprintf("%v", DeviceID),
			City:        city,
			Item:        GetRandItem(city),
			Gender:      gender,
			Order:       order,
			OS:          iphoneOS,
			SourcePage:  sourcePage,
			ProductCity: productCity,
			IsRemote:    isRemote,
			CreateTime:  StartTime.Format("2006-01-02 15:04:05"),
			Hour:        StartTime.Format("15"),
		}
		result = append(result, tmp)
	}

	return result
}

// GenerateFakeGuy 生成异常常的日志
func GenerateFakeGuy() []CSVItem {
	// 总条数
	num := GetRandInt(5)

	userID := GetRandInt(99999999)
	IP := GetRandIpaddr4()
	DeviceID := GetRandDeviceID()
	city := GetRandCity(true)
	productCity := city
	isRemote := GetRandRemote(true)
	order := GetRandOrder()
	iphoneOS := GetRandOS(true)
	sourcePage := GetRandSourcePage()
	gender := GetRandGender()

	if isRemote == "外地" {
		productCity = GetRandCity(true, city)
	}

	StartTime := time.Date(2023, 02, 1+GetRandInt(3), GetRandHour(true), GetRandInt(59), 0, 0, time.Local)

	result := make([]CSVItem, 0, num)
	for i := 0; i < num; i++ {
		randSec := GetRandInt(3)
		StartTime = StartTime.Add(time.Second * time.Duration(randSec))

		tmp := CSVItem{
			ID:          userID,
			IP:          IP,
			DeviceID:    fmt.Sprintf("%v", DeviceID),
			City:        city,
			Item:        GetRandItem(city),
			Gender:      gender,
			Order:       order,
			OS:          iphoneOS,
			SourcePage:  sourcePage,
			ProductCity: productCity,
			IsRemote:    isRemote,
			CreateTime:  StartTime.Format("2006-01-02 15:04:05"),
			Hour:        StartTime.Format("15"),
			Interval:    fmt.Sprintf("%v", randSec),
		}
		result = append(result, tmp)
	}

	return result
}

// GenerateFakeSummary 生成异常汇总日志
func GenerateFakeSummary() CSVItemSummary {
	isRisk := GetRandFake()

	userID := GetRandInt(99999999)
	IP := GetRandIpaddr4()
	DeviceID := GetRandDeviceID()
	city := GetRandCity(false)
	productCity := city
	isRemote := GetRandRemote()
	order := GetRandOrder()
	iphoneOS := GetRandOS(false)
	sourcePage := GetRandSourcePage()
	gender := GetRandGender()

	if isRemote == "外地" {
		productCity = GetRandCity(false, city)
	}

	StartTime := time.Date(2023, 02, 1+GetRandInt(3), GetRandHour(true), GetRandInt(59), 0, 0, time.Local)

	result := CSVItemSummary{
		ID:          userID,
		IP:          IP,
		DeviceID:    fmt.Sprintf("%v", DeviceID),
		City:        city,
		Item:        GetRandItem(city),
		Gender:      gender,
		Order:       order,
		OS:          iphoneOS,
		SourcePage:  sourcePage,
		ProductCity: productCity,
		IsRemote:    isRemote,
		CreateTime:  StartTime.Format("2006-01-02 15:04:05"),
		Hour:        StartTime.Format("15"),
		IsRisk:      isRisk,            // 是否是风险数据
		RiskType:    GetRandFakeType(), // 风险类别
	}
	result.IPCnt = GetRandInt(2)

	// 正常数据
	if isRisk == 0 {
		result.TotalCnt = 1 + GetRandInt(15)
		result.ItemCnt = GetRandInt(5)
		if result.ItemCnt > result.TotalCnt {
			result.ItemCnt = result.TotalCnt
		}
		timeCnt := 0
		for i := 0; i < result.TotalCnt; i++ {
			randSec := GetRandInt(30)
			timeCnt = timeCnt + randSec
		}
		result.TotalTime = timeCnt
		cl := math.Ceil(float64(result.TotalTime) / float64(result.TotalCnt))
		result.AVGTime = int(cl)

		return result
	}

	result.IsRemote = GetRandRemote(true)

	switch result.RiskType {
	case 1: // 点击频次高
		result.TotalCnt = 5 + GetRandInt(15)
		result.IPCnt = GetRandInt(2)
		result.ItemCnt = GetRandInt(5)
		if result.ItemCnt > result.TotalCnt {
			result.ItemCnt = result.TotalCnt
		}
		timeCnt := 0
		for i := 0; i < result.TotalCnt; i++ {
			randSec := GetRandInt(3)
			timeCnt = timeCnt + randSec
		}
		result.TotalTime = timeCnt
		cl := math.Ceil(float64(result.TotalTime) / float64(result.TotalCnt))
		result.AVGTime = int(cl)
	case 2: // ip次数多
		result.TotalCnt = 1 + GetRandInt(15)
		result.IPCnt = result.TotalCnt

		result.ItemCnt = GetRandInt(5)
		if result.ItemCnt > result.TotalCnt {
			result.ItemCnt = result.TotalCnt
		}
		timeCnt := 0
		for i := 0; i < result.TotalCnt; i++ {
			randSec := GetRandInt(30)
			timeCnt = timeCnt + randSec
		}
		result.TotalTime = timeCnt
		cl := math.Ceil(float64(result.TotalTime) / float64(result.TotalCnt))
		result.AVGTime = int(cl)
	}

	return result
}
