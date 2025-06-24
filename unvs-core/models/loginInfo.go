package models

import (
	"dbx"
	"time"
)

type LoginInfo struct {
	dbx.EntityModel
	Id           string  `db:"varchar(36);pk"`
	UserId       string  `db:"varchar(36);idx"` // user id
	RefreshToken string  `db:"varchar(255);idx"`
	AccessToken  *string `db:"varchar(500);idx"`

	CreatedAt time.Time
}

func init() {
	dbx.AddEntities(&LoginInfo{})
}
