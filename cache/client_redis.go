package cache

import (
	"github.com/astaxie/beego/logs"
	"github.com/go-redis/redis/v8"
	"sync"
	"time"
)

/*
*
 使用 redis 做驱动内核
*/

type RedisClient struct {
	client        *redis.Client
	ExpireMap     sync.Map
	DefaultExpire time.Duration
}

//var redisClient *RedisClient

// LoadRedisClient 获取 Redis 相关的缓存器
func LoadRedisClient(c *redis.Client, defaultExpire ...time.Duration) CommonDrive {
	client := &RedisClient{}
	client.client = c

	// 设置默认时长
	if len(defaultExpire) > 0 {
		client.DefaultExpire = defaultExpire[0]
	} else {
		client.DefaultExpire = time.Second * 3600
	}

	return client
}

// GetExpire 获取配置的缓存时长
func (c *RedisClient) GetExpire(k string) time.Duration {
	exp, ok := c.ExpireMap.Load(k)
	if ok && exp != nil {
		return exp.(time.Duration)
	}

	// 如果不设置就默认缓存 1小时
	return time.Second * 3600
}

// SetExpire 设置单个函数的缓存时长
func (c *RedisClient) SetExpire(k string, expire time.Duration) CommonDrive {
	c.ExpireMap.Store(k, expire)
	return c
}

func (c *RedisClient) Save(k string, v interface{}, expire time.Duration) error {
	logs.Info("save")
	b, err := Encode(v)
	logs.Info("saveErr: %v", err)
	if err != nil {
		return err
	}
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

// DeleteFunc 删除对应的函数结果缓存
func (c *RedisClient) DeleteFunc(funcName interface{}, params ...interface{}) bool {
	logs.Info("DeleteFunc")
	cb := &CallFuncBody{
		FuncName: funcName,
		Params:   params,
		cache:    c,
	}

	cacheKey, _, err := cb.GetCacheKey()
	logs.Trace("cacheKey: %v, err: %v", cacheKey, err)
	if err != nil {
		return false
	}

	return c.Delete(cacheKey)
}

// GetCacheFuncKey 返回对应的缓存 key
func (c *RedisClient) GetCacheFuncKey(funcName interface{}, params ...interface{}) (cacheKey string, err error) {
	logs.Info("GetCacheFuncKey")
	cb := &CallFuncBody{
		FuncName: funcName,
		Params:   params,
		cache:    c,
	}

	cacheKey, _, err = cb.GetCacheKey()
	return cacheKey, err
}

// CacheFunc 这里主要做的几件事, 1. 根据方法名和传参获取缓存, 注册返回值类型 key, 2. 查询是否存在对应 key 的缓存结果, 3. 返回缓存/返回执行结果
func (c *RedisClient) CacheFunc(funcName interface{}, params ...interface{}) *CallFuncBody {
	logs.Info("CacheFunc")
	cb := &CallFuncBody{
		FuncName: funcName,
		Params:   params,
		cache:    c,
	}

	cacheKey, funcCacheKey, err := cb.GetCacheKey()
	logs.Trace("cacheKey: %v, err: %v", cacheKey, err)
	if err != nil {
		cb.Err = err
		return cb
	}

	cachedRes := make([]interface{}, 0)
	err = cb.cache.Get(cacheKey, &cachedRes)

	logs.Trace("getCacheResult: %#+v, err: %v", cachedRes, err)
	if err != nil {
		goto STEP1
	}

	cb.Result = cachedRes
	return cb

STEP1:
	res, err := CallFunc(*cb)
	logs.Info("notCache result: %v", res)
	if err != nil {
		cb.Err = err
		return cb
	}

	// 添加
	cb.Result = res

	// 生成新的缓存
	err = cb.cache.Save(cacheKey, res, c.GetExpire(funcCacheKey))
	if err != nil {
		cb.Err = err
		return cb
	}

	return cb
}
