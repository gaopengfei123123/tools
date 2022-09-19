package cache

// CommonDrive 统一缓存驱动接口
type CommonDrive interface {
	Save(k string, v interface{}) error
	Get(k string, target interface{}) error
	Delete(k string) bool
	Exist(k string) bool
	CacheFunc(funcName interface{}, params ...interface{}) *CallFuncBody
}
