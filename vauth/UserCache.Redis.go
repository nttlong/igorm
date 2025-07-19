package unvsauth

import (
	"context"
	"encoding/json"
	"time"
	"unvs-auth/models"

	"github.com/go-redis/redis/v8"
)

type RedisUserCache struct {
	client *redis.Client
	ttl    time.Duration
}

func (r *RedisUserCache) Get(id string) (*models.User, bool) {
	val, err := r.client.Get(context.Background(), "user:"+id).Result()
	if err != nil {
		return nil, false
	}
	var u models.User
	json.Unmarshal([]byte(val), &u)
	return &u, true
}

func (r *RedisUserCache) Set(u *models.User) {
	data, _ := json.Marshal(u)
	r.client.Set(context.Background(), "user:"+u.Email, data, r.ttl)
	r.client.Set(context.Background(), "user:"+u.Username, data, r.ttl)
}
