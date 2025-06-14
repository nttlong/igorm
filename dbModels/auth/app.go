package dbmodels

import "dbx"

type AppConfig struct {
	dbx.EntityModel `db:"AppConfig"`
	Id              int    `db:"pk;auto" json:"id"`
	Name            string `db:"uk;varchar(255)" json:"name"`
	Tenant          string `db:"uk;varchar(255)" json:"tenant"`
	AppId           string `db:"uk;varchar(36)" json:"appId"`
	JwtSecret       string `db:"varchar(500) json:"-"`
}

func init() {
	dbx.AddEntities(&AppConfig{})
}
