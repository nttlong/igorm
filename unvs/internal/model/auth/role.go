package auth

import (
	"dbx"
	"time"
)

type Role struct {
	Id        int        `db:"pk;auto" json:"id"`
	RoleId    string     `db:"varchar(36);uk" json:"roleId"`
	Code      string     `db:"uk;varchar(50)" json:"code"`
	Name      string     `db:"uk;varchar(50)" json:"name"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	CreatedBy string     `db:"varchar(50);idx" json:"createdBy,omitempty"`
	UpdatedBy *string    `db:"varchar(50);idx" json:"updatedBy,omitempty"`

	Description dbx.FullTextSearchColumn `json:"description" swag:"-"`
	Users       []User                   `db:"fk:RoleId" json:"users,omitempty"`
}

func init() {
	dbx.AddEntities(&User{})
}
