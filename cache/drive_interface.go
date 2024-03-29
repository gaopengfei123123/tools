package cache

import (
	"github.com/go-redis/redis/v8"
	"time"
)

// CommonDrive 统一缓存驱动接口
type CommonDrive interface {
	Save(k string, v interface{}, expire time.Duration) error
	Get(k string, target interface{}) error
	Delete(k string) bool
	Exist(k string) bool
	CacheFunc(funcName interface{}, params ...interface{}) *CallFuncBody
	DeleteFunc(funcName interface{}, params ...interface{}) bool
	SetExpire(k string, exp time.Duration) CommonDrive
	GetExpire(k string) time.Duration
	GetCacheFuncKey(funcName interface{}, params ...interface{}) (cacheKey string, err error)
	GetRedisClient() *redis.Client
}
