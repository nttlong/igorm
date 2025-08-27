package models

import (
	"vdb"
)

type User struct {
	vdb.Model[User]
	ID       uint32 `db:"primaryKey;auto"`
	Username string `db:"unique;size:50;"`
	Password string
	Email    *string `db:"unique;size:50;"`
	Phone    *string `db:"unique;size:50;"`
	Active   bool    `db:"default:true"`
}

func init() {
	vdb.ModelRegistry.Add(User{})
}
