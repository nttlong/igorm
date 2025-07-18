package models

import (
	"dbv"
	"time"
)

type User struct {
	dbv.Model[User]

	ID           int    `db:"pk;auto"`
	UserId       string `db:"size:36;unique"`
	Email        string `db:"uk:uq_email;size:150"`
	Phone        string `db:"size:20"`
	Username     string `db:"size:50;unique"`
	HashPassword string `db:"size:100"`

	IsActive        bool `db:"default:true"` // có thể dùng để khóa tài khoản
	EmailVerifiedAt *time.Time
	PhoneVerifiedAt *time.Time

	LastLoginAt      *time.Time
	LoginFailedCount int `db:"default:0"`
	LockedUntil      *time.Time

	Provider   string `db:"size:50"` // oauth2/google/github/...
	ProviderID string `db:"size:100"`

	ResetToken       string `db:"size:100"`
	ResetTokenExpiry *time.Time

	RememberToken string `db:"size:100"`

	BaseModel
}

func init() {
	dbv.ModelRegistry.Add(&User{})
}
