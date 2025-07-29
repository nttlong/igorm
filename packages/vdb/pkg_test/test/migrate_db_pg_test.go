package test

// import (
// 	"vdb"
// 	"vdb/tenantDB"
// 	"testing"

// 	_ "vdb/pkg_test/models"

// 	"github.com/stretchr/testify/assert"
// )

// func TestPGGenerateSQLCreateTable(t *testing.T) {
// 	pgDsn := "postgres://postgres:123456@localhost:5432/a001?sslmode=disable"
// 	// create new migrate instance
// 	db, err := tenantDB.Open("postgres", pgDsn)

// 	assert.NoError(t, err)

// 	migrator, err := vdb.NewMigrator(db)
// 	assert.NoError(t, err)
// 	err = migrator.DoMigrates()
// 	assert.NoError(t, err)

// 	// fmt.Println(sql)

// }
