package vdb

import (
	"reflect"
	"sync"
	"vdb/migrate"
)

type getRepoInfoByTypeInit struct {
	once sync.Once
	val  *getRepoInfoByTypeInfo
}
type getRepoInfoByTypeInfo struct {
	tableName string
	entity    migrate.Entity
}

var cacheGetRepoInfoByType sync.Map

func getRepoInfoByType(typ reflect.Type) *getRepoInfoByTypeInfo {
	actual, _ := cacheGetRepoInfoByType.LoadOrStore(typ, &getRepoInfoByTypeInit{})
	init := actual.(*getRepoInfoByTypeInit)
	init.once.Do(func() {
		model := migrate.GetModelByType(typ)
		if model == nil {
			init.val = nil
			return
		}
		init.val = &getRepoInfoByTypeInfo{
			tableName: model.GetTableName(),
			entity:    model.GetEntity(),
		}
	})

	return init.val
}
func Repo[T any]() *Repository[T] {
	repoType := getRepoInfoByType(reflect.TypeFor[T]())
	if repoType == nil {
		panic(NewModelError(reflect.TypeFor[T]()))
	}

	return &Repository[T]{
		tableName: repoType.tableName,
		entity:    &repoType.entity,
	}
}
func Ptr[T any](t T) *T {
	return &t
}
