package models

import "eorm"

type Contract struct {
	eorm.Model[Contract]
	BaseModel
	ID        int    `db:"pk;auto"`
	UserID    int    `db:"idx:idx_contract_user"`
	StartDate string `db:"type:date"`
	EndDate   string `db:"type:date"`
	Type      string `db:"size:50"` // probation, permanent, seasonal...
	Note      string `db:"size:255"`
}
