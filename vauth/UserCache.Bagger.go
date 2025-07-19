package unvsauth

import (
	"encoding/json"
	"time"
	"unvs-auth/models"

	"github.com/dgraph-io/badger/v4"
)

type BadgerUserCache struct {
	db  *badger.DB
	ttl time.Duration
}

func NewBadgerUserCache(path string, ttl time.Duration) (*BadgerUserCache, error) {
	opts := badger.DefaultOptions(path)
	opts.Logger = nil // optional: disable log
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	return &BadgerUserCache{db: db, ttl: ttl}, nil
}

func (c *BadgerUserCache) Get(id string) (*models.User, bool) {
	var user *models.User
	err := c.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("user:" + id))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			var u models.User
			if err := json.Unmarshal(val, &u); err != nil {
				return err
			}
			user = &u
			return nil
		})
	})
	if err != nil {
		return nil, false
	}
	return user, true
}

func (c *BadgerUserCache) Set(u *models.User) {
	_ = c.db.Update(func(txn *badger.Txn) error {
		data, _ := json.Marshal(u)
		e := badger.NewEntry([]byte("user:"+u.Email), data).WithTTL(c.ttl)
		if err := txn.SetEntry(e); err != nil {
			return err
		}
		e2 := badger.NewEntry([]byte("user:"+u.Username), data).WithTTL(c.ttl)
		return txn.SetEntry(e2)
	})
}

func (c *BadgerUserCache) Delete(id string) {
	_ = c.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte("user:" + id))
	})
}
