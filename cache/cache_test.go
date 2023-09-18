package cache

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/gaopengfei123123/tools"
	"github.com/go-redis/redis/v8"
	"reflect"
	"testing"
	"time"
)

func TestDecode(t *testing.T) {
	info := map[string]string{
		"name":    "C语言中文网",
		"website": "http://c.biancheng.net/golang/",
	}

	bt, err := Encode(info)
	t.Logf("bytes: %v err: %v", bt, err)

	newMp := map[string]string{}
	err = Decode(bt, &newMp)
	t.Logf("v: %v,  err: %v", newMp, err)
}

func TestDecode2(t *testing.T) {
	info := map[string]string{
		"name":    "GPF",
		"website": "https://blog.justwe.site",
	}

	bt, err := Encode(info)
	t.Logf("bytes: %v err: %v", bt, err)

	newMp := map[string]int{}
	err = Decode(bt, &newMp)
	t.Logf("v: %v,  err: %v", newMp, err)
}

func getRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:         "127.0.0.1:6379",
		Password:     "",
		MinIdleConns: 5,
		PoolSize:     20,
	})
}

func TestLoadRedisClient(t *testing.T) {
	redisClient := getRedisClient()
	key1 := "TestLoadRedisClient"
	err := LoadRedisClient(redisClient).Save(key1, "xxx", time.Second*10)
	t.Logf("err: %v", err)

	info := map[string]string{
		"name":    "GPF",
		"website": "https://blog.justwe.site",
	}
	err = LoadRedisClient(redisClient).Save("TestLoadRedisClientInfo", info, time.Second*10)
	t.Logf("method:Save err: %v", err)

	// 从缓存中获取
	newInfo := map[string]string{}
	err = LoadRedisClient(redisClient).Get("TestLoadRedisClientInfo", &newInfo)
	t.Logf("method:Get    info: %#+v, err: %v", newInfo, err)

	res := LoadRedisClient(redisClient).Exist(key1)
	t.Logf("method:Exist key: %v exist: %v", key1, res)

	res = LoadRedisClient(redisClient).Delete(key1)
	t.Logf("method:Delete delete key;%v res: %v", key1, res)

	res = LoadRedisClient(redisClient).Exist(key1)
	t.Logf("method:Exist  key: %v exist: %v", key1, res)
}

func TestRedisClient_Save(t *testing.T) {
	redisClient := getRedisClient()
	var tOffset int
	logs.Trace("resxxx: %#+v", tOffset)
	tmpConfig := TopicConfig{
		//"topic1": map[int32]int64{
		//	12: 34,
		//},
		//"topic2": map[int32]int64{
		//	22: 33,
		//},
	}
	tmpConfig = SetTopicConfig(tmpConfig, "topic1", 12, 15)
	tmpConfig = SetTopicConfig(tmpConfig, "topic2", 12, 15)
	logs.Trace(tmpConfig)
	v, ok := GetTopicConfig(tmpConfig, "topic1", 12)
	logs.Trace("getTopicConfig: %v %v", v, ok)
	v, ok = GetTopicConfig(tmpConfig, "topicxx", 12)
	logs.Trace("getTopicConfig: %v %v", v, ok)
	err := LoadRedisClient(redisClient).Save("TestLoadRedisClientInfoXXX", tmpConfig, time.Second*300)
	logs.Trace("save err: %v", err)
	//var resConfig TopicConfig
	//err = LoadRedisClient(redisClient).Get("TestLoadRedisClientInfoXXX", &resConfig)
	//logs.Trace("resConfig: %#+v", resConfig)

	res := TopicConfig{}
	LoadRedisClient(redisClient).Get("TestLoadRedisClientInfoXXX", &res)
	logs.Trace("result: %v", res)
}

func SetTopicConfig(conf TopicConfig, topic string, partition int32, offset int64) TopicConfig {
	_, exist := conf[topic]
	if !exist {
		conf[topic] = make(map[int32]int64)
	}
	conf[topic][partition] = offset
	return conf
}

func GetTopicConfig(conf TopicConfig, topic string, partition int32) (int64, bool) {
	part, exist := conf[topic]
	if !exist {
		return 0, false
	}
	offset, exist := part[partition]
	return offset, exist
}

type TopicConfig map[string]map[int32]int64

// 缓存方法示例
func TestRedisClient_CacheFunc(t *testing.T) {
	redisClient := getRedisClient()
	sign := "cache_message2"

	// 初始化缓存工具
	cache := LoadRedisClient(redisClient)
	// 设置 Demo函数结果缓存时间
	cache.SetExpire(GetFuncName(Demo), time.Second*180)

	for i := 0; i < 5; i++ {
		var errMsg error
		var funcRes string
		err := cache.CacheFunc(Demo, sign).GetResult(&funcRes, &errMsg)
		t.Logf("method: CacheFunc err: %v, result: %v, funcErr: %v", err, funcRes, errMsg)
		time.Sleep(time.Second)
	}
}

// 注册多个缓存方法, 每个函数缓存的时长不一样
func TestCallFuncBody_GetResult(t *testing.T) {
	cache := LoadRedisClient(getRedisClient())
	cache.SetExpire(GetFuncName(Demo), time.Second*180)
	cache.SetExpire(GetFuncName(Demo2), time.Second*3600)

	var resultMsg string
	cache.CacheFunc(Demo, "params xxx").GetResult(&resultMsg)
	logs.Trace("Func: Demo,  expire: %v, result: %v", cache.GetExpire(GetFuncName(Demo)), resultMsg)

	// 必须提前注册类型, 不然解不出来
	gob.Register(&map[string]string{})
	// TODO 这个有问题, map, 切片, 对象等传址的参数不能很好的识别
	var mp map[string]string
	tmp := map[string]string{
		"Name": "GPFXX",
	}
	cache.CacheFunc(Demo2, tmp).GetResult(&mp)
	logs.Trace("Func: Demo2,   expire: %v, result: %v", cache.GetExpire(GetFuncName(Demo2)), mp)
}

// 缓存复杂类型的数据测试
func TestCallFuncBody_GetResult2(t *testing.T) {
	cache := LoadRedisClient(getRedisClient())

	var result *DemoResult
	result = &DemoResult{}
	gob.Register(result)
	err := cache.CacheFunc(Demo3, int32(3)).GetResult(&result)
	logs.Trace("Func demo3, err: %v, result: %#+v", err, result)

	tools.PrintJson("cache result", result)

	result = Demo3(3)
	logs.Trace("NoCache: %#+v", result)
}

func TestRedisClient_DeleteFunc(t *testing.T) {
	cache := LoadRedisClient(getRedisClient())
	res := cache.DeleteFunc(Demo3, int32(3))
	logs.Trace("delete res: %v", res)
}

func TestRedisClient_CacheFunc2(t *testing.T) {
	m0 := reflect.TypeOf(Demo2)

	for i := 0; i < m0.NumOut(); i++ {
		logs.Trace("返回值: %#+v", m0.Out(i))
		logs.Trace("索引: %v", m0.Out(i).String())
	}
}

// gob: cannot encode nil pointer of type *errcode.Error inside interface
// 这里调整了 cache/cache.go:18 的逻辑
func TestRedisClient_CacheFunc3(t *testing.T) {
	cache := LoadRedisClient(getRedisClient())
	result := cache.CacheFunc(Demo4, 3)
	logs.Trace("res: %#+v", result)

	var res int
	var err *Error
	result.GetResult(&res, &err)
	logs.Trace("result: %#+v, %#+v", res, err)
}

func TestRedisClient_GetCacheFuncKey(t *testing.T) {
	cache := LoadRedisClient(getRedisClient())

	result := cache.CacheFunc(Demo4, 3)
	logs.Trace("res: %#+v", result)

	key, err := cache.GetCacheFuncKey(Demo4, 3)
	logs.Trace("key: %v, err: %v", key, err)

	cache.SetExpire(key, time.Second*5)
	exp := cache.GetExpire(key)
	logs.Trace("exp: %v", exp.Seconds())
}

func TestEncode(t *testing.T) {
	var dao bytes.Buffer

	tmp := map[string]string{
		"Name": "GPF",
	}
	encoder := gob.NewEncoder(&dao)

	err := encoder.Encode(tmp)
	logs.Trace("err: %v", err)
	logs.Trace("encode: %v", dao.String())
}

// 用来测试缓存的函数1
func Demo(msg string) (string, error) {
	t := time.Now()
	tt := t.Format("06-01-02 15:04:05")
	res := fmt.Sprintf("message: %s, time: %s", msg, tt)
	return res, nil
}

// 用来测试缓存的函数2
func Demo2(mp map[string]string) (map[string]string, error) {
	t := time.Now()
	tt := t.Format("06-01-02 15:04:05")
	res := fmt.Sprintf("map: %s, time: %s", mp, tt)
	mp["timer"] = res
	return mp, nil
}

func Demo4(id int) (int, *Error) {
	return id, nil
}

type Error struct {
	code           int
	httpStatusCode int
	msg            string
	details        []string
}

type DemoResult struct {
	ID          int32
	ListType    []DemoResultItem
	ListTypePtr []*DemoResultItem
	Name        string
	Child       DemoResultItem
	ChildPtr    *DemoResultItem
	Date        string
}

type DemoResultItem struct {
	Key   string
	Value string
}

func Demo3(id int32) *DemoResult {
	result := &DemoResult{}
	result.ID = id

	demoList1 := make([]DemoResultItem, 0)
	result.ListType = append(demoList1, DemoResultItem{"For", "Bar"}, DemoResultItem{"For2", "Bar"})

	demoListPtr := make([]*DemoResultItem, 0)
	result.ListTypePtr = append(demoListPtr, &DemoResultItem{"For2", "Bar2"}, &DemoResultItem{"For22", "Bar2"})

	result.Name = "xxx"
	result.Child = DemoResultItem{"For3", "Bar3"}
	result.ChildPtr = &DemoResultItem{"For4", "Bar4"}

	result.Date = time.Now().Format("2006-01-02 15:04:05")

	return result
}
