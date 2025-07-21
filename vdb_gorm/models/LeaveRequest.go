package models

import "time"

// LeaveRequest struct GORM
type LeaveRequest struct {
	ID         int `gorm:"primaryKey;autoIncrement"` // vdb: db:"pk;auto"
	EmployeeId int `gorm:"index"`                    // vdb: db:"idx"
	StartDate  time.Time
	EndDate    time.Time
	Reason     string `gorm:"type:varchar(255)"` // vdb: db:"size:255"
	Status     string `gorm:"type:varchar(20)"`  // vdb: db:"size:20"
	BaseModel         // Nhúng BaseModel
	// Quan hệ với Employee (cho GORM)
	Employee Employee `gorm:"foreignKey:EmployeeId"`
}

// Tương tự, không cần init() function cho khóa ngoại trong GORM.
