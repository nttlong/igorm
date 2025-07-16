package dbv

import (
	"dbv/migrate"
	"reflect"
	"sync"
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
		init.val = &getRepoInfoByTypeInfo{
			tableName: migrate.GetModelByType(typ).GetTableName(),
			entity:    migrate.GetModelByType(typ).GetEntity(),
		}
	})

	return init.val
}
func Repo[T any]() *Repository[T] {
	repoType := getRepoInfoByType(reflect.TypeFor[T]())

	return &Repository[T]{
		tableName: repoType.tableName,
		entity:    &repoType.entity,
	}
}
func Ptr[T any](t T) *T {
	return &t
}
