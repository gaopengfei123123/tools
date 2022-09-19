package cache

import (
	"fmt"
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
		Addr:         "",
		Password:     "",
		MinIdleConns: 5,
		PoolSize:     20,
	})
}

func TestLoadRedisClient(t *testing.T) {
	redisClient := getRedisClient()
	key1 := "TestLoadRedisClient"
	err := LoadRedisClient(redisClient).Save(key1, "xxx")
	t.Logf("err: %v", err)

	info := map[string]string{
		"name":    "GPF",
		"website": "https://blog.justwe.site",
	}
	err = LoadRedisClient(redisClient).Save("TestLoadRedisClientInfo", info)
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

func TestRedisClient_CacheFunc(t *testing.T) {
	res, err := Demo("no cache message")
	t.Logf("no cache msg: %v, err: %v", res, err)

	redisClient := getRedisClient()
	sign := "cache_message"

	var errMsg error
	var funcRes string
	err = LoadRedisClient(redisClient).CacheFunc(Demo, sign).GetResult(&funcRes, &errMsg)
	t.Logf("method: CacheFunc err: %v", err)
}

func Demo(msg string) (string, error) {
	t := time.Now()
	tt := t.Format("06-01-02 15:04:05")
	res := fmt.Sprintf("message: %s, time: %s", msg, tt)
	return res, nil
}
