package models

import "vdb"

type Department struct {
	vdb.Model[Department]

	ID       int    `db:"pk;auto"`
	Name     string `db:"size:100;uk:uq_dept_name"`
	Code     string `db:"size:20;uk:uq_dept_code"`
	ParentID *int
	BaseModel
}

func init() {
	(&Department{}).AddForeignKey("ParentID", &Department{}, "ID", &vdb.CascadeOption{
		OnDelete: false,
		OnUpdate: false,
	})

}
