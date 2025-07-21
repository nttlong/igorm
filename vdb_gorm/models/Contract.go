package models

import "time"

// Contract struct GORM
type Contract struct {
	ID        int `gorm:"primaryKey;autoIncrement"` // vdb: db:"pk;auto"
	UserID    int `gorm:"index:idx_contract_user"`  // vdb: db:"idx:idx_contract_user"
	StartDate time.Time
	EndDate   time.Time
	Type      string `gorm:"type:varchar(50)"`  // vdb: db:"size:50"
	Note      string `gorm:"type:varchar(255)"` // vdb: db:"size:255"
	BaseModel        // Nhúng BaseModel
	// Quan hệ với User (cho GORM)
	User User `gorm:"foreignKey:UserID"`
}
