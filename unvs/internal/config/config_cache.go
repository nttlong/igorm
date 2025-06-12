package config

import (
	cache "caching"
	"reflect"
	"sync"
	"time"
)

var cacher cache.Cache
var onceCache sync.Once

func GetCache() cache.Cache {
	onceCache.Do(func() {
		cacher = cache.NewInMemoryCache(
			reflect.TypeOf(Config{}),
			10*time.Minute, 10*time.Minute)
	})
	return cacher

}
