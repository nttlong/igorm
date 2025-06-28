package unvsef

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createPgDb() *sql.DB {
	dsn := "user=postgres password=123456 host=localhost port=5432 dbname=fx001 sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	return db
}
func TestConnect(t *testing.T) {
	db := createPgDb()
	assert.NotNil(t, db)
	defer db.Close()
}
func TestGetSchemaPostgres(t *testing.T) {
	db := createPgDb()
	assert.NotNil(t, db)
	defer db.Close()

	d := NewPostgresDialect(db)
	schema, err := d.GetSchema(db, "fx001")
	assert.Nil(t, err)
	assert.Equal(t, "public", schema)
}
func TestEntityPostgres(t *testing.T) {

}
