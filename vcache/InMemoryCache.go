package vcache

import (
	"context"
	"log"
	"reflect"
	"time"

	gocache "github.com/patrickmn/go-cache"
)

type InMemoryCache struct {
	client    *gocache.Cache
	prefixKey string
}

// NewInMemoryCache tạo một instance mới của InMemoryCache
// use "github.com/patrickmn/go-cache"
func NewInMemoryCache(

	defaultExpiration,
	cleanupInterval time.Duration) Cache {

	//strHasKey is string version of hashKey

	// for in-memory cache, default expiration and cleanup interval are ignored
	return &InMemoryCache{
		client: gocache.New(defaultExpiration, cleanupInterval),
	} // no check error here, so just return nil
}

// Get implements Cache.Get for InMemoryCache
func (c *InMemoryCache) Get(ctx context.Context, key string, dest interface{}) bool {

	val := reflect.ValueOf(dest)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	} else {
		log.Println("InMemoryCache: Can not get object, dest muste be a pointer")
	}
	destType := reflect.TypeOf(dest).Elem()
	key = c.prefixKey + ":" + key + ":" + destType.PkgPath() + "." + destType.Name()

	r, f := c.client.Get(key)

	if !f {
		return false
	}
	// dest = r
	val.Set(reflect.ValueOf(r))
	// if err := bytesDecodeObject(r.([]byte), val.Interface()); err != nil {
	// 	log.Println("InMemoryCache: Không thể decode object: ", err)
	// 	return false
	// }

	return true
}

// Set implements Cache.Set for InMemoryCache
func (c *InMemoryCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) {
	val := reflect.ValueOf(value)
	destType := reflect.TypeOf(value)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		destType = destType.Elem()
	}

	key = c.prefixKey + ":" + key + ":" + destType.PkgPath() + "." + destType.Name()
	// objBff, err := bytesEncodeObject(val.Interface())
	// if err != nil {
	// 	log.Println("InMemoryCache: Không thể encode object: ", err)
	// 	return
	// }

	if ttl == 0 { // Sử dụng TTL mặc định nếu được truyền 0
		c.client.Set(key, val.Interface(), gocache.DefaultExpiration)
	} else {
		c.client.Set(key, val.Interface(), ttl)
	}
}

// Delete implements Cache.Delete for InMemoryCache
func (c *InMemoryCache) Delete(ctx context.Context, key string) {

	c.client.Delete(c.prefixKey + ":" + key)
}

// Close implements Cache.Close for InMemoryCache (không cần làm gì cho go-cache)
func (c *InMemoryCache) Close() error {
	log.Println("InMemoryCache đã đóng.")
	return nil
}
