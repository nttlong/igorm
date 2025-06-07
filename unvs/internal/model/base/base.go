package base

import (
	"dbx"
	_ "dbx"
	"time"
)

type BaseModel struct {
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	CreatedBy string     `db:"varchar(50);idx" json:"createdBy,omitempty"`
	UpdatedBy *string    `db:"varchar(50);idx" json:"updatedBy,omitempty"`

	Description dbx.FullTextSearchColumn `json:"description" swag:"-"`
}
