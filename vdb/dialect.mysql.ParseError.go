package vdb

import (
	"fmt"

	"github.com/go-sql-driver/mysql"
)

func (d *mySqlDialect) ParseError(err error) error {
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		if mysqlErr.Number == 1048 {
			return d.ParseError1048(mysqlErr)

		}
		fmt.Println(mysqlErr.Number)

	}
	return err
}
