package models

import (
	"time"
	"vdb"
)

type SecurityPolicy struct {
	vdb.Model[SecurityPolicy]
	ID               int        `db:"pk;auto"`
	TenantID         string     `db:"size:50;uk"` // Ràng buộc duy nhất theo Tenant
	MaxLoginFailures int        `db:"default:5"`
	LockoutMinutes   int        `db:"default:15"`
	JwtSecret        string     `db:"size:255"`
	JwtExpireMinutes int        `db:"default:60"`
	CreatedAt        time.Time  `db:"default:now"`
	UpdatedAt        *time.Time `db:"default:now"`
}

func init() {
	vdb.ModelRegistry.Add(&SecurityPolicy{})
}
