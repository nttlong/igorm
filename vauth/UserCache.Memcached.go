package vauth

import (
	"encoding/json"
	"vauth/models"

	"github.com/bradfitz/gomemcache/memcache"
)

type MemcachedUserCache struct {
	client *memcache.Client
	ttl    int32 // in seconds
}

func NewMemcachedUserCache(addr string, ttl int32) *MemcachedUserCache {
	return &MemcachedUserCache{
		client: memcache.New(addr),
		ttl:    ttl,
	}
}

func (c *MemcachedUserCache) Get(id string) (*models.User, bool) {
	item, err := c.client.Get("user:" + id)
	if err != nil {
		return nil, false
	}
	var user models.User
	if err := json.Unmarshal(item.Value, &user); err != nil {
		return nil, false
	}
	return &user, true
}

func (c *MemcachedUserCache) Set(u *models.User) {
	data, _ := json.Marshal(u)
	_ = c.client.Set(&memcache.Item{Key: "user:" + u.Email, Value: data, Expiration: c.ttl})
	_ = c.client.Set(&memcache.Item{Key: "user:" + u.Username, Value: data, Expiration: c.ttl})
}

func (c *MemcachedUserCache) Delete(id string) {
	_ = c.client.Delete("user:" + id)
}
