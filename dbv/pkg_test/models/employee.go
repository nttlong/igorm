package models

import (
	"dbv"
)

type Employee struct {
	dbv.Model[Employee]
	ID           int    `json:"id" db:"pk;auto;"`
	FirstName    string `json:"name" db:"size(50);idx"`
	LastName     string `json:"lastName" db:"size(50);idx"`
	DepartmentID int    `json:"departmentId" db:"fk(Department.ID)"`
	PositionID   int    `json:"positionId" db:"fk(Position.ID)"`
	UserID       int    `json:"userId" db:"fk(User.ID)"`
	BaseModel
}

func init() {
	dbv.ModelRegistry.Add(Employee{}, &Department{})
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
