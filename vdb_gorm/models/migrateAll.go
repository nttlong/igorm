package models

import "gorm.io/gorm"

func MigrateAllModels(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&BaseModel{}, // BaseModel thường không cần migrate riêng nếu nó chỉ là nhúng
		&Department{},
		&Position{},
		&Employee{},
		&Attendance{},
		&Contract{},
		&LeaveRequest{},
		&Salary{},
	)
}
