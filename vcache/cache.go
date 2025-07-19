// internal/app/cache/cache.go
package vcache

import (
	"context"
	"time"
)

/*
this interface provide basic cache operations
In order to
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
	Close() error                           // Đóng kết nối/giải phóng tài nguyên của cache
}

// === Triển khai InMemoryCache sử dụng github.com/patrickmn/go-cache ===

// InMemoryCache là triển khai của Cache interface sử dụng go-cache
