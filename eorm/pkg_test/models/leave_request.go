package models

import (
	"dbv"
	"time"
)

type LeaveRequest struct {
	dbv.Model[LeaveRequest]
	BaseModel
	ID        int `db:"pk;auto"`
	UserID    int `db:"idx:idx_leave_user"`
	StartDate time.Time
	EndDate   time.Time
	Reason    string `db:"size:255"`
	Status    string `db:"size:20"` // pending, approved, rejected
}

func init() {
	dbv.ModelRegistry.Add(&LeaveRequest{})
	(&LeaveRequest{}).AddForeignKey("UserID", &User{}, "ID", nil)

}
