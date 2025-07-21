package models

import "time"

// BaseModel GORM
type BaseModel struct {
	// GORM tự động quản lý CreatedAt và UpdatedAt nếu có tên này và kiểu time.Time
	// Thẻ `column` để đảm bảo tên cột là snake_case trong DB.
	// Thẻ `index` để tạo index.
	CreatedAt *time.Time `gorm:"column:created_at;index"`
	UpdatedAt *time.Time `gorm:"column:updated_at;index"`

	Description *string `gorm:"type:varchar(255)"` // vdb: db:"size:255"
}
