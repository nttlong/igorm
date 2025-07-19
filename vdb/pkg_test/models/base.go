package models

import "time"

type BaseModel struct {
	CreatedAt   *time.Time `db:"default:now;idx;default:now"`
	UpdatedAt   *time.Time `db:"default:now;idx"`
	Description *string    `db:"size:255"`
}
