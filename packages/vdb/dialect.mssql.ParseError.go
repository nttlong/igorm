package vdb

import (
	"vdb/migrate"

	mssql "github.com/microsoft/go-mssqldb"
)

func (d *mssqlDialect) ParseError(dbSchame *migrate.DbSchema, err error) error {
	//go-mssqldb.Error
	if mssqlErr, ok := err.(mssql.Error); ok {
		return d.Error2627(mssqlErr)
	}

	return err
}
