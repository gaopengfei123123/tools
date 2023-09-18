package mock

import (
	"github.com/astaxie/beego/logs"
	"testing"
)

func TestGetRandInt(t *testing.T) {
	for i := 0; i < 100; i++ {
		t.Logf("rand: %v \n", GetRandInt(100))
	}
}

func TestRandPool_LoadConfig(t *testing.T) {
	mp := map[string]int{
		"北京": 100,
		"上海": 80,
		"杭州": 50,
		"广东": 50,
		"深圳": 45,
		"重庆": 40,
		"成都": 30,
	}

	cp := new(RandPool)
	cp.LoadConfig(mp)
	logs.Trace("cityPool: %v", cp)
}

func TestRandPool_GetItem(t *testing.T) {
	mp := map[string]int{
		"北京": 100,
		"上海": 80,
		"杭州": 50,
		"广东": 50,
		"深圳": 45,
		"重庆": 40,
		"成都": 30,
	}

	cp := new(RandPool)
	cp.LoadConfig(mp)

	result := make(map[string]int)

	for i := 0; i < 10000; i++ {
		v := cp.GetItem()
		result[v] += 1
	}
	logs.Trace("result: %v", result)
}

func TestGetRandCity(t *testing.T) {
	result := make(map[string]int)

	for i := 0; i < 10000; i++ {
		v := GetRandCity(false)
		result[v] += 1
	}
	logs.Trace("result: %v", result)
}

func TestGetRandItem(t *testing.T) {
	result := make(map[string]int)

	for i := 0; i < 10000; i++ {
		v := GetRandItem("xxx")
		result[v] += 1
	}
	logs.Trace("result: %v", result)
}

func TestGetRandGender(t *testing.T) {
	result := make(map[string]int)

	for i := 0; i < 10000; i++ {
		v := GetRandGender()
		result[v] += 1
	}
	logs.Trace("result: %v", result)
}

func TestGetRandOrder(t *testing.T) {
	result := make(map[string]int)

	for i := 0; i < 10000; i++ {
		v := GetRandOrder()
		result[v] += 1
	}
	logs.Trace("result: %v", result)
}

func TestGetRandOS(t *testing.T) {
	result := make(map[string]int)

	for i := 0; i < 10000; i++ {
		v := GetRandOS(false)
		result[v] += 1
	}
	logs.Trace("result: %v", result)
}

func TestGetRandSourcePage(t *testing.T) {
	result := make(map[string]int)

	for i := 0; i < 10000; i++ {
		v := GetRandSourcePage()
		result[v] += 1
	}
	logs.Trace("result: %v", result)
}

func TestGetRandRemote(t *testing.T) {
	result := make(map[string]int)

	for i := 0; i < 10000; i++ {
		v := GetRandRemote()
		result[v] += 1
	}
	logs.Trace("result: %v", result)
}

func TestGetRandHour(t *testing.T) {
	result := make(map[int]int)

	for i := 0; i < 10000; i++ {
		v := GetRandHour(false)
		result[v] += 1
	}
	logs.Trace("result: %v", result)
}

func TestMockDataGenerator(t *testing.T) {
	MockDataGenerator()
}

// 最终生成原始数据文档
func TestOutputCsv(t *testing.T) {
	//OutputCsvDemo()
	OutPutCsv(300)
}

// 最终生成汇总数据文档
func TestOutPutCsvSummary(t *testing.T) {
	OutPutCsvSummary(1000)
}

func TestGetRandIpaddr4(t *testing.T) {
	logs.Trace(GetRandIpaddr4())
}

func TestGetRandDeviceID(t *testing.T) {
	logs.Trace(GetRandDeviceID())
}

func TestGenerateNormalGuy(t *testing.T) {
	logs.Trace("%v", GenerateNormalGuy())
}

func TestGetRandFake(t *testing.T) {
	result := make(map[int]int)

	for i := 0; i < 10000; i++ {
		v := GetRandFake()
		result[v] += 1
	}
	logs.Trace("result: %v", result)
}
