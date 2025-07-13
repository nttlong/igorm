package models

import "eorm"

type LeaveRequest struct {
	eorm.Model[LeaveRequest]
	BaseModel
	ID        int    `db:"pk;auto"`
	UserID    int    `db:"idx:idx_leave_user"`
	StartDate string `db:"type:date"`
	EndDate   string `db:"type:date"`
	Reason    string `db:"size:255"`
	Status    string `db:"size:20"` // pending, approved, rejected
}
