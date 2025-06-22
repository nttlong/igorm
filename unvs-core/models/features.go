package models

import (
	"dbx"
	"time"
)

type Features struct {
	dbx.EntityModel
	Id              string `db:"varchar(36);pk" json:"id"`
	Name            string `db:"varchar(255);uk" json:"name"`
	Description     dbx.FullTextSearchColumn
	CreatedAt       time.Time         `db:"idx" json:"created_at"`
	UpdatedAt       *time.Time        `db:"idx" json:"updated_at"`
	CreatedBy       string            `db:"varchar(255);idx" json:"created_by"`
	UpdatedBy       *string           `db:"varchar(255);idx" json:"updated_by"`
	FeaturesDetails []FeaturesDetails `db:"fk:FeaturesId"`
}
type FeaturesDetails struct {
	dbx.EntityModel
	Id          int        `db:"auto;pk" json:"id"`
	FeaturesId  string     `db:"varchar(36)" json:"features_id"`
	Module      string     `db:"varchar(255);idx" json:"module"`
	Action      string     `db:"varchar(255);idx" json:"action"`
	Description string     `db:"varchar(255)" json:"description"`
	CreatedAt   time.Time  `db:"idx" json:"created_at"`
	UpdatedAt   *time.Time `db:"idx" json:"updated_at"`
	CreatedBy   string     `db:"varchar(255);idx" json:"created_by"`
}

func init() {
	dbx.AddEntities(&Features{}, &FeaturesDetails{})
}
