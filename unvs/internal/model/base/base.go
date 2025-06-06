package base

import (
	"dbx"
	_ "dbx"
	"time"
)

type BaseModel struct {
	CreatedAt time.Time
	UpdatedAt *time.Time
	CreatedBy string  `db:"varchar(50);idx"`
	UpdatedBy *string `db:"varchar(50);idx"`

	Description dbx.FullTextSearchColumn
}
