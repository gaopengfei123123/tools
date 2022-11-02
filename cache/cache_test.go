package cache

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/go-redis/redis/v8"
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
	var tOffset int
	logs.Info("resxxx: %#+v", tOffset)
	tmpConfig := TopicConfig{
		//"topic1": map[int32]int64{
		//	12: 34,
		//},
		//"topic2": map[int32]int64{
		//	22: 33,
		//},
	}
	tmpConfig = SetTopicConfig(tmpConfig, "topic1", 12, 15)
	logs.Info(tmpConfig)
	//err = LoadRedisClient(redisClient).Save("TestLoadRedisClientInfoXXX", tmpConfig, time.Second*10)
	//var resConfig TopicConfig
	//err = LoadRedisClient(redisClient).Get("TestLoadRedisClientInfoXXX", &resConfig)
	//logs.Info("resConfig: %#+v", resConfig)

}

func SetTopicConfig(conf TopicConfig, topic string, partition int32, offset int64) TopicConfig {
	_, exist := conf[topic]
	if !exist {
		conf = make(map[string]map[int32]int64)
	}
	_, exist = conf[topic][partition]
	if !exist {
		conf[topic] = make(map[int32]int64)
	}
	conf[topic][partition] = offset
	return conf
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

	gob.Register(&map[string]string{})

	//var resultMsg string
	//cache.CacheFunc(Demo, "params xxx").GetResult(&resultMsg)
	//logs.Info("Func: Demo,  expire: %v, result: %v", cache.GetExpire(GetFuncName(Demo)), resultMsg)

	// TODO 这个有问题, map, 切片, 对象等传址的参数不能很好的识别
	var mp map[string]string
	tmp := map[string]string{
		"Name": "GPF",
	}
	cache.CacheFunc(Demo2, tmp).GetResult(&mp)
	logs.Info("Func: Demo,  expire: %v, result: %v", cache.GetExpire(GetFuncName(Demo2)), mp)
}

func TestEncode(t *testing.T) {
	var dao bytes.Buffer

	tmp := map[string]string{
		"Name": "GPF",
	}
	encoder := gob.NewEncoder(&dao)

	err := encoder.Encode(tmp)
	logs.Info("err: %v", err)
	logs.Info("encode: %v", dao.String())
}

func Demo(msg string) (string, error) {
	t := time.Now()
	tt := t.Format("06-01-02 15:04:05")
	res := fmt.Sprintf("message: %s, time: %s", msg, tt)
	return res, nil
}

func Demo2(mp map[string]string) (map[string]string, error) {
	t := time.Now()
	tt := t.Format("06-01-02 15:04:05")
	res := fmt.Sprintf("map: %s, time: %s", mp, tt)
	mp["timer"] = res
	return mp, nil
}
