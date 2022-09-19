package cache

import (
	"github.com/astaxie/beego/logs"
	"github.com/go-redis/redis/v8"
	"time"
)

/*
*
 使用 redis 做驱动内核
*/

type RedisClient struct {
	client *redis.Client
}

//var redisClient *RedisClient

// LoadRedisClient 获取 Redis 相关的缓存器
func LoadRedisClient(c *redis.Client) CommonDrive {
	client := &RedisClient{}
	client.client = c
	return client
}

func (c *RedisClient) GetExpire(k string) time.Duration {
	return time.Second * 180
}

func (c *RedisClient) Save(k string, v interface{}) error {
	logs.Info("save")
	b, err := Encode(v)
	if err != nil {
		return err
	}
	expire := c.GetExpire(k)
	err = c.client.Set(c.client.Context(), k, string(b), expire).Err()
	return err
}

func (c *RedisClient) Get(k string, target interface{}) error {
	logs.Info("Get")
	val, err := c.client.Get(c.client.Context(), k).Result()
	if err != nil {
		return err
	}
	err = Decode([]byte(val), target)
	return err
}

func (c *RedisClient) Delete(k string) bool {
	logs.Info("Delete: %v", k)
	i, err := c.client.Del(c.client.Context(), k).Result()
	if err != nil || i <= 0 {
		return false
	}
	return true
}

func (c *RedisClient) Exist(k string) bool {
	logs.Info("Exist %v", k)
	i, err := c.client.Exists(c.client.Context(), k).Result()
	if err != nil || i <= 0 {
		return false
	}
	return true
}

func (c *RedisClient) CacheFunc(funcName interface{}, params ...interface{}) *CallFuncBody {
	logs.Info("CacheFunc")
	bd := &CallFuncBody{
		FuncName: funcName,
		Params:   params,
		cache:    c,
	}
	return bd
}
