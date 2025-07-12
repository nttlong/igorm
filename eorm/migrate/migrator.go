package migrate

import (
	"reflect"
)

type IMigrator interface {
	/*
		this function should return a map of go type to database column type mapping
		for example:
		map[reflect.Type]string{
			reflect.TypeOf(int(0)): "int", //<-- int in SQL server type
			reflect.String: "navarchar", //<-- nvarchar in SQL server type
			reflect.TypeOf(time.Time{}): "datetime2", //<-- datetime2 in SQL server type
		or
		map[reflect.Type]string{
			reflect.TypeOf(int(0)): "integer", //<-- integer in MySQL type
			reflect.String: "varchar", //<-- varchar in MySQL type
			reflect.TypeOf(time.Time{}): "datetime", //<-- datetime in MySQL type

	*/
	GetColumnDataTypeMapping() map[reflect.Type]string

	/*
		this function should return a SQL create table statement for the given go type
		for example:

		type User struct {
		  		eorm.Model `db:"table:users"` //<-- if tag table:users found, it will be used as table name
											  // else use name of struct the pluralized it for table name
		  		ID        int    `db:"pk"` //<--  field name is name of Go struct field, pluralized it for column name if Database
										  // else use name of struct the singularized it for primary key column name

				Name      string `db:"length:255"` //<-- for length tag, it will be used as length of column in Database.
													// Heed: just use if GO field type is string

				Email     string `db:"unique"` //<-- for unique tag, it will be used as unique constraint in Database.
			}
		sql make table looks like
			CREATE TABLE users (
				id int PRIMARY KEY,
				name nvarchar(255),
				email nvarchar(255)
			)


	*/
	GetSqlCreateTable(reflect.Type) (string, error)
}
