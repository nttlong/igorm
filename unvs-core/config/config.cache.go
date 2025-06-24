package config

import (
	"reflect"
	"sync"
	"time"

	cacher "unvs.core/cacher"
)

var cache cacher.Cache
var onceCache sync.Once

func getCache() cacher.Cache {
	var _cacher cacher.Cache

	if AppConfigInstance.Cache == "bagger" {
		cacheBagger, err := cacher.NewBadgerCache(
			reflect.TypeOf(Config{}),
			"unvs",
		)
		if err != nil {
			panic(err)
		}
		_cacher = cacheBagger

	}
	if AppConfigInstance.Cache == "in-memory" {
		_cacher = cacher.NewInMemoryCache(
			reflect.TypeOf(Config{}),
			10*time.Minute, 10*time.Minute)

	}
	if AppConfigInstance.Cache == "memcached" {
		_cacher = cacher.NewMemcachedCache(
			reflect.TypeOf(Config{}),
			AppConfigInstance.MemcachedServer,
		)

	}
	if AppConfigInstance.Cache == "redis" {
		_cacher = cacher.NewRedisCache(
			reflect.TypeOf(Config{}),
			AppConfigInstance.Redis.Host,
			AppConfigInstance.Redis.Password,
			AppConfigInstance.Redis.DB,
			time.Duration(AppConfigInstance.Redis.Timeout)*time.Millisecond,
		)

	}

	return _cacher
}
func GetCache() cacher.Cache {
	onceCache.Do(func() {
		cache = getCache()
	})
	return cache

}
