package models

import (
	"time"
	"vdb"
)

type Role struct {
	vdb.Model[Role]
	ID          int       `db:"pk;auto"`
	TenantID    string    `db:"size:36;uk:uk_role_name_tenant"`  // duy nhất trên tenant
	Name        string    `db:"size:100;uk:uk_role_name_tenant"` // duy nhất theo tenant
	Description *string   `db:"size:255"`
	CreatedAt   time.Time `db:"default:now()"`
	UpdatedAt   *time.Time
}

type UserRole struct {
	vdb.Model[UserRole]
	ID        int       `db:"pk;auto"`
	AccountID int       `db:"idx:idx_userrole_account_role"` // khóa ngoại tới account
	RoleID    int       `db:"idx:idx_userrole_account_role"` // khóa ngoại tới role
	TenantID  string    `db:"size:36;idx"`                   // để xác định theo tenant
	CreatedAt time.Time `db:"default:now()"`
}

func init() {
	vdb.ModelRegistry.Add(&Role{}, &UserRole{})
	(&UserRole{}).AddForeignKey("AccountID", &Account{}, "ID", nil)
	(&UserRole{}).AddForeignKey("RoleID", &Role{}, "ID", nil)
}
