package vdb

import (
	"fmt"
	"reflect"
	"sync"
	"vdb/migrate"
	"vdb/tenantDB"
)

type Model[T any] struct {
	migrate.Entity
}
type initModelMakeAlias struct {
	once sync.Once
	val  string
}

var cacheModelMakeAlias = sync.Map{}

func modelMakeAlias(typ reflect.Type, alias string) string {
	key := fmt.Sprintf("%s:%s", typ.String(), alias)
	actual, _ := cacheModelMakeAlias.LoadOrStore(key, &initModelMakeAlias{})
	init := actual.(*initModelMakeAlias)
	init.once.Do(func() {
		repoType := inserterObj.getEntityInfo(typ)
		init.val = repoType.tableName + " AS " + alias
	})
	return init.val
}

func (m *Model[T]) As(alias string) string {

	return modelMakeAlias(reflect.TypeFor[T](), alias)
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
