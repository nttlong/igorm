package models

import (
	"time"
	"vdb"
)

type LeaveRequest struct {
	vdb.Model[LeaveRequest]
	BaseModel
	ID         int `db:"pk;auto"`
	EmployeeId int `db:"idx"`
	StartDate  time.Time
	EndDate    time.Time
	Reason     string `db:"size:255"`
	Status     string `db:"size:20"` // pending, approved, rejected
}

func init() {
	vdb.ModelRegistry.Add(&LeaveRequest{})
	(&LeaveRequest{}).AddForeignKey("EmployeeId", &Employee{}, "ID", nil)

}
