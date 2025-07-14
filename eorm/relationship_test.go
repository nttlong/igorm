package eorm

import "testing"

type Department struct {
	Model[Department]

	ID       int    `db:"pk;auto"`
	Name     string `db:"size:100;uk:uq_dept_name"`
	Code     string `db:"size:20;uk:uq_dept_code"`
	ParentId *int
}

func (d *Department) Build() {
	d.AddForeignKey("ParentId", &Department{}, "ID")

}
func init() {
	ModelRegistry.Add(&Department{})
}
func TestRelationship(t *testing.T) {
	pk := (&Department{ID: 1}).AddForeignKey("ParentId", &Department{}, "ID")
	t.Log(pk)

}
