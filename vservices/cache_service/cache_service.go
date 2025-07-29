package cache_service

// import (
// 	"fmt"

// 	"vcache"
// )

// func NewCacheService(cfg *config.ConfigService) (vcache.Cache, error) {
// 	c := (*cfg).Get()

// 	switch c.CacheType {
// 	case "redis":
// 		redis := c.Redis
// 		cache := vcache.NewRedisCache(
// 			redis.Nodes, redis.Password, redis.PrefixKey,
// 			redis.DB, redis.Timeout,
// 		)
// 		return cache, nil

// 	case "memcached":
// 		return vcache.NewMemcachedCache(c.Memcached.Nodes, c.Memcached.PrefixKey), nil

// 	case "badger":
// 		cache, err := vcache.NewBadgerCache(c.Badger.Path, c.Badger.PrefixKey)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return cache, nil

// 	case "inmemory":
// 		cache := vcache.NewInMemoryCache(
// 			c.InMemory.DefaultTTL,
// 			c.InMemory.CleanupInterval,
// 		)
// 		return cache, nil

// 	default:
// 		return nil, fmt.Errorf("unsupported cache type: %s", c.CacheType)
// 	}
// }
