package test

import (
	_ "eorm"
	"testing"
)

func TestQueryBuilderPG(t *testing.T) {
	pgDsn := "postgres://postgres:postgres@localhost:5432/testdb?sslmode=disable"
	t.Log(pgDsn)
}
