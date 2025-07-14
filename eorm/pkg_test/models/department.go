package models

import "eorm"

type Department struct {
	eorm.Model[Department]
	BaseModel
	ID       int    `db:"pk;auto"`
	Name     string `db:"size:100;uk:uq_dept_name"`
	Code     string `db:"size:20;uk:uq_dept_code"`
	ParentID *int
}

func init() {
	(&Department{}).AddForeignKey("ParentID", &Department{}, "ID")

}
