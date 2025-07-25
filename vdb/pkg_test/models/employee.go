package models

import (
	"vdb"
)

type Employee struct {
	vdb.Model[Employee]
	ID           int    `json:"id" db:"pk;auto;"`
	FirstName    string `json:"name" db:"size(50);idx"`
	LastName     string `json:"lastName" db:"size(50);idx"`
	DepartmentID int    `json:"departmentId"`
	PositionID   int    `json:"positionId"`
	UserID       int    `json:"userId"`
	BaseModel
}

func init() {
	vdb.ModelRegistry.Add(Employee{}, &Department{})
	(&Employee{}).AddForeignKey(
		"DepartmentID",
		&Department{},
		"ID", nil).AddForeignKey(
		"PositionID",
		&Position{},
		"ID", nil).AddForeignKey(
		"UserID",
		&User{},
		"ID", nil,
	)

}
