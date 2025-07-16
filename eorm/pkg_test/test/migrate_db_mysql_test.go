package test

import (
	"dbv"
	"dbv/tenantDB"
	"testing"

	_ "dbv/pkg_test/models"

	"github.com/stretchr/testify/assert"
)

func TestMySqlGenerateSQLCreateTable(t *testing.T) {

	mySqlDsn := "root:123456@tcp(127.0.0.1:3306)/a001"
	// create new migrate instance
	db, err := tenantDB.Open("mysql", mySqlDsn)

	assert.NoError(t, err)

	migrator, err := dbv.NewMigrator(db)

	assert.NoError(t, err)
	err = migrator.DoMigrates()
	assert.NoError(t, err)

	// fmt.Println(sql)

}
