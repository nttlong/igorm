// internal/app/cache/cache.go
package vcache

import (
	"context"
	"time"
)

/*
this interface provide basic cache operations
In order to
Usage:

	1- Create a new instance with InMemoryCache
		var cache vcache.Cache
		cache = vcache.NewInMemoryCache(10*time.Second, 10*time.Second)
	2- Create a new Cache with Bagger:
		var cache vcache.Cache
		cache, err := vcache.NewBadgerCache(<path to db>,  <prefix key>)
	3- Create a new Cache with Redis:
	   var cache vcache.Cache
		cache = vcache.NewRedisCache(<server>, <password>, <prefix key>, 0, 10*time.Second)
	4- Create a new Cache with Memcached:
		var cache vcache.Cache
		cache = vcache.NewMemcachedCache([server1, server2], <prefix key>)
	Heed: all cache implementations were already tested and proven to work correctly.
*/
type Cache interface {

	// get object from cache
	// example: Get("key", &obj)
	// @description: This function will combine  key and package path of object and name of object type to create a unique key for cache.
	// @param key: string, key of object in cache, actually it is a part of real cache key
	Get(ctx context.Context, key string, dest interface{}) bool // Lấy giá trị từ cache

	// @description: This function will combine  key and package path of object and name of object type to create a unique key for cache.
	// @param key: string, key of object in cache, actually it is a part of real cache key
	// @param value: interface{}, value of object to store in cache
	// @param ttl: time.Duration, time to live of object in cache 0 is default value which means no expiration
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) // Đặt giá trị vào cache với TTL
	// Đặt giá trị vào cache với TTL
	Delete(ctx context.Context, key string) // Xóa một key khỏi cache
	Close() error
	GetBool(ctx context.Context, key string) (bool, bool)
	// Đóng kết nối/giải phóng tài nguyên của cache
	// Expire(ctx context.Context, key string, ttl time.Duration) // Đặt thời gian hết hạn cho một key)
}

// === Triển khai InMemoryCache sử dụng github.com/patrickmn/go-cache ===

// InMemoryCache là triển khai của Cache interface sử dụng go-cache
