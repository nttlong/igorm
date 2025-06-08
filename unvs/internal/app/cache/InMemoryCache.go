package cache

import (
	"log"
	"time"

	gocache "github.com/patrickmn/go-cache"
)

type InMemoryCache struct {
	client *gocache.Cache
}

// NewInMemoryCache tạo một instance mới của InMemoryCache
func NewInMemoryCache(defaultExpiration, cleanupInterval time.Duration) *InMemoryCache {
	return &InMemoryCache{
		client: gocache.New(defaultExpiration, cleanupInterval),
	}
}

// Get implements Cache.Get for InMemoryCache
func (c *InMemoryCache) Get(key string) (interface{}, bool) {
	return c.client.Get(key)
}

// Set implements Cache.Set for InMemoryCache
func (c *InMemoryCache) Set(key string, value interface{}, ttl time.Duration) {
	if ttl == 0 { // Sử dụng TTL mặc định nếu được truyền 0
		c.client.Set(key, value, gocache.DefaultExpiration)
	} else {
		c.client.Set(key, value, ttl)
	}
}

// Delete implements Cache.Delete for InMemoryCache
func (c *InMemoryCache) Delete(key string) {
	c.client.Delete(key)
}

// Close implements Cache.Close for InMemoryCache (không cần làm gì cho go-cache)
func (c *InMemoryCache) Close() error {
	log.Println("InMemoryCache đã đóng.")
	return nil
}
