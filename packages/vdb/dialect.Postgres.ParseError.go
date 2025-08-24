package vdb

import (
	"errors"
	"fmt"

	"vdb/migrate"

	"github.com/lib/pq"
)

func (d *postgresDialect) ParseError(dbSchema *migrate.DbSchema, err error) error {
	var pgErr *pq.Error
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return d.ParseError23505(dbSchema, pgErr)

		}
		panic(fmt.Errorf(`not implemented,vdb\dialect.Postgres.go`))
	} else {
		return err
	}

}
