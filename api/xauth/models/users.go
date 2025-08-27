package models

import (
	"time"
	"vdb"
)

type User struct {
	vdb.Model[User]
	ID        uint32 `db:"primaryKey;auto"`
	Username  string `db:"unique;size:50;"`
	Password  string
	Email     *string   `db:"unique;size:50;"`
	Phone     *string   `db:"unique;size:50;"`
	Active    bool      `db:"default:true"`
	CreatedOn time.Time `db:"default:now()"`
	UpdatedOn *time.Time
}

func init() {
	vdb.ModelRegistry.Add(User{})
}
