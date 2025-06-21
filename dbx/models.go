package dbx

import "time"

type Tenants struct {
	EntityModel
	Id          string     `db:"varchar(36);pk" json:"id"`
	Name        string     `db:"varchar(255);uk" json:"name"`
	DbName      string     `db:"varchar(255);uk" json:"dbName"`
	Description string     `db:"varchar(255)" json:"description"`
	CreatedAt   time.Time  `db:"idx" json:"createdAt"`
	Updated     *time.Time `db:"idx" json:"updatedAt"`
	CreatedBy   string     `db:"varchar(36)" json:"createdBy"`
	UpdatedBy   *string    `db:"varchar(36)" json:"updatedBy"`
}
