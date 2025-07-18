package models

import (
	"dbv"
	"time"
)

type User struct {
	dbv.Model[User]
	BaseModel
	ID         int    `db:"pk;auto"`
	Name       string `db:"size:100"`
	Email      string `db:"uk:uq_email;size:150"`
	Gender     string `db:"size:10"` // male, female, other
	Birthday   time.Time
	Phone      string  `db:"size:20"`
	Address    string  `db:"size:255"`
	DeptID     int     `db:"idx:idx_user_dept"`
	PositionID int     `db:"idx:idx_user_pos"`
	Username   *string `db:"size:50;unique"`
}
