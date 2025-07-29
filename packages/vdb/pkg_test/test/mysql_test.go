package test

import (
	"fmt"
	"testing"
	"vdb"
	"vdb/tenantDB"

	"vdb/pkg_test/models"
	_ "vdb/pkg_test/models"

	"github.com/stretchr/testify/assert"
)

type HrmRepo struct {
	users *vdb.Repository[models.User]
}

func TestSelectAll(t *testing.T) {
	mySqlDsn := "root:123456@tcp(127.0.0.1:3306)/a001?multiStatements=true"
	// create new migrate instance
	db, err := tenantDB.Open("mysql", mySqlDsn)

	assert.NoError(t, err)

	migrator, err := vdb.NewMigrator(db)
	assert.NoError(t, err)
	err = migrator.DoMigrates()
	if err != nil {
		fmt.Println(err)
	}
	rows, err := vdb.SelectAll[models.User](db)
	assert.NoError(t, err)
	assert.Greater(t, len(rows), 151073) //<--lay 151.073 dong

}
