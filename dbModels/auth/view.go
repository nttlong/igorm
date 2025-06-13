package dbmodels

import (
	"dbx"
	"time"
)

type View struct {
	dbx.EntityModel `db:"View"`
	Id              int    `db:"pk;auto" json:"id"`
	ViewId          string `db:"varchar(36);uk" json:"viewId"`
	Name            string `db:"varchar(255)"  json:"name"`

	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	CreatedBy   string     `db:"varchar(255)" json:"createdBy"`
	UpdatedBy   *string    `db:"varchar(255)" json:"updatedBy"`
	Description dbx.FullTextSearchColumn
	ViewRoles   []ViewRole `db:"fk:ViewId" json:"viewRoles"`
}
type ViewRole struct {
	dbx.EntityModel `db:"ViewRole"`
	Id              int `db:"pk;auto" json:"id"`
	ViewId          int `db:uk:ViewId_RoleId" json:"viewId"`
	RoleId          int `db:uk:ViewId_RoleId" json:"roleId"`
}

func init() {
	dbx.AddEntities(&View{}, &ViewRole{})
}
