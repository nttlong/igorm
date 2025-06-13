package dbx

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

func mhysqlExecCreateTable(db *sql.DB, dbname string, key string, sqlList SqlCommandList) error {

	if db == nil {
		return fmt.Errorf("please open db first")
	}
	for _, sqlCmd := range sqlList {
		fmt.Println(red+"EXEC: "+reset+sqlCmd.String(), red+"Error: "+reset)
		_, err := db.Exec(sqlCmd.String())
		if err != nil {

			if mySQlErr, ok := err.(*mysql.MySQLError); ok {
				if mySQlErr.Number == 1060 || mySQlErr.Number == 1061 || mySQlErr.Number == 1826 {

					continue
				} else {

					fmt.Println(red+"SQL: "+reset+sqlCmd.String(), red+"Error: "+reset+err.Error())
					return DBXMigrationError{
						Message: err.Error(),
						DBName:  dbname,
						Code:    fmt.Sprintf("%d", mySQlErr.Number), // mysql error code
						Sql:     sqlCmd.String(),
						Err:     err,
					}
				}

			} else {
				fmt.Println(red+"SQL: "+reset+sqlCmd.String(), red+"Error: "+reset+err.Error())

				return DBXMigrationError{
					Message: err.Error(),
					DBName:  dbname,
					Code:    "unknown", // mysql error code
					Sql:     sqlCmd.String(),
					Err:     err,
				}
			}

		}

	}
	//save entityType to cache
	checkCreateTable.Store(key, true)
	return nil
}
