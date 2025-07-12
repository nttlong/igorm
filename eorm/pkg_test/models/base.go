package models

import "time"

type BaseModel struct {
	CreatedAt time.Time `db:"default:now;type:datetime"`
	UpdatedAt time.Time `db:"default:now;type:datetime"`
}
