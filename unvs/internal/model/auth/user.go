package auth

import (
	"unvs/internal/model/base"
)

type User struct {
	Id           int    `db:"pk;auto"`
	UserId       string `db:"varchar(36);uk"`
	Username     string `db:"uk;varchar(255)"`
	PasswordHash string `db:"varchar(255)"`
	Email        string `db:"uk;varchar(320)"`
	base.BaseModel
}
