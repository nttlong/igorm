package models

import (
	"time"
	"vdb"

	"github.com/google/uuid"
)

type User struct {
	vdb.Model[User]
	ID        uint32    `db:"primaryKey;auto"`
	UserId    uuid.UUID `db:"unique;size:36;default:uuid()"`
	Username  string    `db:"unique;size:50;"`
	Password  string    `db:"size:300;"`
	Email     *string   `db:"unique;size:50;"`
	Phone     *string   `db:"unique;size:50;"`
	Active    bool      `db:"default:true"`
	CreatedOn time.Time `db:"default:now()"`
	UpdatedOn *time.Time
}

func init() {
	vdb.ModelRegistry.Add(User{})
}
