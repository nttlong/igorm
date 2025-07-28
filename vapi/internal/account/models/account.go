package models

import (
	"time"
	"vdb"
)

type Account struct {
	vdb.Model[Account]
	ID     int    `db:"pk;auto"`
	UserID string `db:"size:36;uk:uk_account_user_id_tenant"` // duy nhất trên tenant

	Username       string `db:"size:100;uk:uk_account_username_tenant"` // duy nhất trên tenant
	HashedPassword string `db:"size:255"`                               // đã hash
	IsLocked       bool
	LockedUntil    *time.Time
	FullName       string    `db:"size:255"` // thông tin hiển thị
	Email          *string   `db:"size:255;"`
	Phone          *string   `db:"size:20;"`
	Role           string    `db:"size:50;default:'user'"` // Có thể mở rộng RBAC
	CreatedAt      time.Time `db:"default:now()"`
	UpdatedAt      *time.Time
}

type AccountProfile struct {
	vdb.Model[AccountProfile]
	ID        int `db:"pk"`
	AccountID int
	FullName  string     `db:"size:255"`
	Gender    string     `db:"size:10;default:'unknown'"` // male, female, other
	Dob       *time.Time // ngày sinh
	AvatarURL *string    `db:"size:500;"`  // đường dẫn avatar
	Bio       *string    `db:"size:1000;"` // mô tả ngắn
}
type AccountSetting struct {
	vdb.Model[AccountSetting]
	ID           int `db:"pk"`
	AccountID    int
	Language     string `db:"size:10;default:'en'"`         // ví dụ: vi, en, ja
	Timezone     string `db:"size:100;default:'UTC'"`       // tên múi giờ IANA
	DateFormat   string `db:"size:50;default:'YYYY-MM-DD'"` // định dạng ngày
	TimeFormat   string `db:"size:20;default:'HH:mm'"`      // định dạng giờ
	NumberFormat string `db:"size:20;default:'###,###.##'"` // định dạng số
	CurrencyCode string `db:"size:10;default:'USD'"`        // mã tiền tệ
}

func init() {

	vdb.ModelRegistry.Add(&Account{}, &AccountProfile{}, &AccountSetting{})
	(&AccountProfile{}).AddForeignKey("ID", &Account{}, "ID", nil)
	(&AccountSetting{}).AddForeignKey("ID", &Account{}, "ID", nil)
}
