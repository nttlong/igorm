package main

import "dbx"

func getMssqlConfig() dbx.Cfg {
	return dbx.Cfg{
		Driver: "mssql",
		Host:   "localhost",
		// Port:     1433,
		User:     "sa",
		Password: "123456",
		SSL:      false,
	}
}
func main() {
	dbx := dbx.NewDBX(getMssqlConfig())
	dbx.Open()
	dbx.Ping()

}
