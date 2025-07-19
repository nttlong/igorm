package models

import (
	"time"
	"vdb"
)

type User struct {
	vdb.Model[User]

	ID           int    `db:"pk;auto"`
	UserId       string `db:"size:36;uk"`
	Email        string `db:"uk;size:150"`
	Phone        string `db:"size:20"`
	Username     string `db:"size:50;uk"`
	HashPassword string `db:"size:100"`

	IsActive        bool
	EmailVerifiedAt *time.Time
	PhoneVerifiedAt *time.Time

	LastLoginAt      *time.Time
	LoginFailedCount int
	LockedUntil      *time.Time

	Provider   string `db:"size:50"` // oauth2/google/github/...
	ProviderID string `db:"size:100"`

	ResetToken       string `db:"size:100"`
	ResetTokenExpiry *time.Time

	RememberToken string `db:"size:100"`

	BaseModel
}

func init() {
	vdb.ModelRegistry.Add(&User{})
}
