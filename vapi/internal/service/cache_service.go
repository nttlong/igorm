package service

import (
	"fmt"
	"vcache"
)

type CacheService interface {
	Get() vcache.Cache
}

type cacheServiceImpl struct {
	cache vcache.Cache
}

func NewCacheService(cfg *ConfigService) (CacheService, error) {
	c := (*cfg).Get()

	switch c.CacheType {
	case "redis":
		redis := c.Redis
		cache := vcache.NewRedisCache(
			redis.Nodes, redis.Password, redis.PrefixKey,
			redis.DB, redis.Timeout,
		)
		return &cacheServiceImpl{cache: cache}, nil

	case "memcached":
		return &cacheServiceImpl{cache: vcache.NewMemcachedCache(c.Memcached.Nodes, c.Memcached.PrefixKey)}, nil

	case "badger":
		cache, err := vcache.NewBadgerCache(c.Badger.Path, c.Badger.PrefixKey)
		if err != nil {
			return nil, err
		}
		return &cacheServiceImpl{cache: cache}, nil

	case "inmemory":
		cache := vcache.NewInMemoryCache(
			c.InMemory.DefaultTTL,
			c.InMemory.CleanupInterval,
		)
		return &cacheServiceImpl{cache: cache}, nil

	default:
		return nil, fmt.Errorf("unsupported cache type: %s", c.CacheType)
	}
}

func (s *cacheServiceImpl) Get() vcache.Cache {
	return s.cache
}
