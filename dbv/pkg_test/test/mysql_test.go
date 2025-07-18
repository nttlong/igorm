package test

import (
	"dbv"
	"dbv/tenantDB"
	"fmt"
	"testing"

	"dbv/pkg_test/models"
	_ "dbv/pkg_test/models"

	"github.com/stretchr/testify/assert"
)

type HrmRepo struct {
	users *dbv.Repository[models.User]
}

func TestSelectAll(t *testing.T) {
	mySqlDsn := "root:123456@tcp(127.0.0.1:3306)/a001?multiStatements=true"
	// create new migrate instance
	db, err := tenantDB.Open("mysql", mySqlDsn)

	assert.NoError(t, err)

	migrator, err := dbv.NewMigrator(db)
	assert.NoError(t, err)
	err = migrator.DoMigrates()
	if err != nil {
		fmt.Println(err)
	}
	rows, err := dbv.SelectAll[models.User](db)
	assert.NoError(t, err)
	assert.Greater(t, len(rows), 151073) //<--lay 151.073 dong

}
