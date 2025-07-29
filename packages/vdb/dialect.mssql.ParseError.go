package vdb

import (
	mssql "github.com/microsoft/go-mssqldb"
)

func (d *mssqlDialect) ParseError(err error) error {
	//go-mssqldb.Error
	if mssqlErr, ok := err.(mssql.Error); ok {
		return d.Error2627(mssqlErr)
	}

	return err
}
