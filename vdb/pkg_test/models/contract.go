package models

import (
	"time"
	"vdb"
)

type Contract struct {
	vdb.Model[Contract]
	BaseModel
	ID        int `db:"pk;auto"`
	UserID    int `db:"idx:idx_contract_user"`
	StartDate time.Time
	EndDate   time.Time
	Type      string `db:"size:50"` // probation, permanent, seasonal...
	Note      string `db:"size:255"`
}
