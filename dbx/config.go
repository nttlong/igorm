package dbx

import "fmt"

/**
Configuration for database connection
/*/
type Cfg struct {
	Driver         string // postgres, mysql, mssql
	Host           string // hostname or ip address
	Port           int    // port number
	User           string // username
	Password       string // password
	SSL            bool   //	use SSL connection
	DbName         string // if is IsMultiTenancy the database is a database manage all tenancies else is a single database
	IsMultiTenancy bool   // is multitenant database
}

/*
Create connection string for different databases
Support for Postgres, MySQL, MSSQL
*/
func (c *Cfg) dns(dbname string) string {

	if c.Driver == "postgres" {

		return c.makeDnsPostgres(dbname)
	} else if c.Driver == "mysql" {
		return c.makeDnsMySql(dbname)
	} else if c.Driver == "mssql" {
		return c.makeDnsMssql(dbname)
	} else {
		panic(fmt.Errorf("unsupported driver %s", c.Driver))
	}

}
