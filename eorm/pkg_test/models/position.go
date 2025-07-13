package models

import "eorm"

type Position struct {
	eorm.Model[Position]
	BaseModel
	ID    int    `db:"pk;auto"`
	Title string `db:"size:100;uk:uq_pos_title"`
	Level int    `db:""`
}
