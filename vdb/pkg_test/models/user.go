package models

import "vdb"

type User struct {
	vdb.Model[User]

	ID     int     `db:"pk;auto"`
	UserId *string `db:"size:36;unique"`

	Email string `db:"uk:uq_email;size:150"`

	Phone string `db:"size:20"`

	Username     *string `db:"size:50;unique"`
	HashPassword *string `db:"size:100"`
	BaseModel
}
