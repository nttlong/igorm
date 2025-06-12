package config

import (
	cache "caching"
	"reflect"
	"sync"
	"time"
)

var cacher cache.Cache
var onceCache sync.Once

func getCache() cache.Cache {
	var _cacher cache.Cache
	if AppConfigInstance.Cache == "bagger" {
		cacher, err := cache.NewBadgerCache(
			reflect.TypeOf(Config{}),
			"unvs",
		)
		if err != nil {
			panic(err)
		}
		_cacher = cacher

	}
	if AppConfigInstance.Cache == "in-memory" {
		cacher = cache.NewInMemoryCache(
			reflect.TypeOf(Config{}),
			10*time.Minute, 10*time.Minute)
	}
	if AppConfigInstance.Cache == "memcached" {
		_cacher = cache.NewMemcachedCache(
			reflect.TypeOf(Config{}),
			AppConfigInstance.MemcachedServer,
		)

	}
	if AppConfigInstance.Cache == "redis" {
		_cacher = cache.NewRedisCache(
			reflect.TypeOf(Config{}),
			AppConfigInstance.Redis.Host,
			AppConfigInstance.Redis.Password,
			AppConfigInstance.Redis.DB,
			time.Duration(AppConfigInstance.Redis.Timeout)*time.Millisecond,
		)

	}

	return _cacher
}
func GetCache() cache.Cache {
	onceCache.Do(func() {
		cacher = getCache()
	})
	return cacher

}
