package eorm

import (
	"eorm/migrate"
	"eorm/tenantDB"
	"reflect"
)

type Model[T any] struct {
	migrate.Entity
	obj interface{}
}

func (m *Model[T]) New() Model[T] {
	obj := reflect.New(reflect.TypeFor[T]()).Interface().(T)
	ret := Model[T]{obj: obj}
	return ret

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
