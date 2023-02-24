package mock

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/gocarina/gocsv"
	"os"
	"time"
)

func MockDataGenerator() {
	logs.Info("xxx")
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
		logs.Info("err: %v", err)
		logs.Info("end")
		return
	}
}

func OutPutCsv(total int) {
	f, err := os.OpenFile("pv_data.csv", os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	//关闭流
	defer f.Close()
	//写入UTF-8 格式
	f.WriteString("\xEF\xBB\xBF")

	res := make([]CSVItem, 0, total*10)
	for i := 0; i < total; i++ {
		list := GenerateNormalGuy()

		res = append(res, list...)
	}

	err = gocsv.Marshal(res, f)

	if err != nil {
		logs.Info("err: %v", err)
		logs.Info("end")
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
}

// GenerateNormalGuy 生成正常的日志
func GenerateNormalGuy() []CSVItem {
	// 总条数
	num := GetRandInt(20)

	userID := GetRandInt(99999999)
	IP := GetRandIpaddr4()
	DeviceID := GetRandDeviceID()
	city := GetRandCity()
	productCity := city
	isRemote := GetRandRemote()
	order := GetRandOrder()
	iphoneOS := GetRandOS()
	sourcePage := GetRandSourcePage()
	gender := GetRandGender()

	if isRemote == "外地" {
		productCity = GetRandCity(city)
	}

	StartTime := time.Date(2023, 02, 1+GetRandInt(3), GetRandHour(), GetRandInt(59), 0, 0, time.Local)

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
