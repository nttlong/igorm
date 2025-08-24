package dbmodels

import (
	"time"
	"vdb"
)

type Users struct {
	vdb.Model[Users]

	// Khóa chính
	ID int64 `db:"pk;auto"`

	// Thông tin đăng nhập
	Username       string  `db:"unique;size:50"`
	Email          *string `db:"unique;size:100"`
	PhoneNumber    *string `db:"unique;size:20"` // nếu muốn xác thực bằng số ĐT
	HashedPassword string  `db:"size:255;idx"`

	// Thông tin cá nhân
	FullName    *string `db:"size:100;idx"`
	DateOfBirth *time.Time
	Gender      string `db:"size:10"` // male, female, other

	// Trạng thái & bảo mật
	IsActive       bool  `db:"default:true"`
	IsLocked       *bool `db:"default:false"`
	FailedAttempts int   `db:"default:0"` // số lần đăng nhập sai
	LastLoginAt    *time.Time
	LastLoginIP    string `db:"size:45;idx"` // IPv6 max length 45 chars

	// Metadata
	CreatedAt time.Time  `db:"autoCreateTime"`
	UpdatedAt *time.Time `db:"autoUpdateTime"`
}

func init() {
	vdb.RegisterModel(&Users{})
}
