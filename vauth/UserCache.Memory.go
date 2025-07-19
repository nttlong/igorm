package unvsauth

import (
	"sync"
	"unvs-auth/models"
)

type MemoryUserCache struct {
	data sync.Map
}

func (c *MemoryUserCache) Get(id string) (*models.User, bool) {
	v, ok := c.data.Load(id)
	if !ok {
		return nil, false
	}
	return v.(*models.User), true
}

func (c *MemoryUserCache) Set(u *models.User) {
	c.data.Store(u.Email, u)
	c.data.Store(u.Username, u)
}
func (c *MemoryUserCache) Delete(id string) {
	c.data.Delete(id)
}
func (c *MemoryUserCache) Invalidate(u *models.User) {
	c.data.Delete(u.Email)
	c.data.Delete(u.Username)
}
