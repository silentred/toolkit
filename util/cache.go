package util

type Cache interface {
	Set(key string, obj interface{}) bool
	Get(key string) (interface{}, bool)
	Del(key string) bool
}

// // MemCache do memory-cache
// type MemCache struct {
// 	c      *cache.Cache
// 	expire time.Duration
// }

// func NewMemCache(expireSecond int) *MemCache {
// 	expire := time.Duration(expireSecond) * time.Second
// 	checkInterval := expire / 4
// 	if checkInterval < time.Second {
// 		checkInterval = time.Second
// 	}
// 	c := cache.New(expire, checkInterval)
// 	return &MemCache{
// 		c:      c,
// 		expire: expire,
// 	}
// }

// func (mc *MemCache) Get(key string) (interface{}, bool) {
// 	return mc.c.Get(key)
// }

// func (mc *MemCache) Set(key string, obj interface{}) bool {
// 	mc.c.Set(key, obj, mc.expire)
// 	return true
// }

// func (mc *MemCache) Del(key string) bool {
// 	mc.c.Delete(key)
// 	return true
// }

type Callable func() interface{}

func TryCache(c Cache, key string, call Callable) interface{} {
	obj, exists := c.Get(key)
	if exists {
		return obj
	}
	obj = call()
	ok := c.Set(key, obj)
	if ok {
		return obj
	}
	return nil
}
