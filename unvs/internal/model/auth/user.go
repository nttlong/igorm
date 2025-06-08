package auth

import (
	"dbx"
	_ "dbx"
	"time"
)

type User struct {
	Id           int        `db:"pk;auto" swag:"-" json:"-"`
	UserId       string     `db:"varchar(36);uk" json:"userId"`
	Username     string     `db:"uk;varchar(255)" json:"username"`
	PasswordHash string     `db:"varchar(255)" json:"-"` // don't expose password hash in API
	Email        string     `db:"uk;varchar(320)" json:"email"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    *time.Time `json:"updatedAt,omitempty"`
	CreatedBy    string     `db:"varchar(50);idx" json:"createdBy,omitempty"`
	UpdatedBy    *string    `db:"varchar(50);idx" json:"updatedBy,omitempty"`

	Description dbx.FullTextSearchColumn `json:"description" swag:"-"`

	IsLocked             bool       `db:"df:false" json:"isLocked"`
	LastPasswordChangeAt *time.Time ` json:"lastPasswordChangeAt,omitempty"`
	LastLoginAt          *time.Time ` json:"lastLoginAt,omitempty"`
	LastFailedLoginAt    *time.Time ` json:"lastFailedLoginAt,omitempty"`
	FailedLoginCount     int        `db:"df:0" json:"failedLoginCount"`
	IsSupperUser         bool       `db:"df:false" json:"isSupperUser"`
	RoleId               *int       `db:"fk:auth_role(Id)" json:"roleId"`
}

func init() {
	dbx.AddEntities(&User{})
}
