package views

import (
	"dbx"
	"time"
)

type View struct {
	ViewId      int64                    `db:"pk;auto" json:"viewId"`
	Path        string                   `db:"uk;varchar(255)" json:"path"`
	Name        string                   `db:"varchar(255)" json:"name"`
	CreatedAt   time.Time                `json:"createdAt"`
	Description dbx.FullTextSearchColumn `json:"description"`
}
type ViewRole struct {
	Id     int64 `db:"pk;auto" json:"id"`
	RoleId int64 `json:"roleId"`
	ViewId int64 `json:"viewId"`
}

func init() {
	dbx.AddEntities(&View{}, &ViewRole{})
}
