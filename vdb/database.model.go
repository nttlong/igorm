package vdb

import (
	"reflect"
	"vdb/migrate"
	"vdb/tenantDB"
)

type Model[T any] struct {
	migrate.Entity
}

func (m *Model[T]) Insert(db *tenantDB.TenantDB) error {

	migrator, err := migrate.NewMigrator(db)
	if err != nil {
		return err
	}
	err = migrator.DoMigrate(reflect.TypeFor[T]())
	if err != nil {
		return err
	}

	return nil
}
