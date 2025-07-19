package vexample

import (
	"testing"
	"vdb"

	"github.com/stretchr/testify/assert"
)

func initDb(driver string, conn string) (*vdb.TenantDB, error) {
	db, err := vdb.Open(driver, conn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
func TestInitDb(t *testing.T) {
	db, err := initDb("mysql", "root:123456@tcp(localhost:3306)/a001")
	assert.NoError(t, err)
	defer db.Close()
}
