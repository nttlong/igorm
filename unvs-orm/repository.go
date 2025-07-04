package orm

import (
	"fmt"
	"reflect"

	internal "unvs-orm/internal"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/microsoft/go-mssqldb"
)

func createErrorRepoFromType(typ reflect.Type, err error) interface{} {
	ret := reflect.New(typ).Elem()
	ret.FieldByName("Err").Set(reflect.ValueOf(err))
	return ret.Interface()
}

func Repository[T any]() T {
	typ := reflect.TypeFor[T]()
	if _, ok := typ.FieldByName("Base"); ok {
		// baseRepo := Base{}

		retVal, err := internal.Utils.GetOrCreateRepository(typ)
		// tenantDbFieldVal := retVal.PtrValueOfRepo.FieldByName("Base")
		// tenantDbFieldVal.Set(reflect.ValueOf(baseRepo))

		if err != nil {
			ret := createErrorRepoFromType(typ, err)
			return ret.(T)
		}

		return retVal.PtrValueOfRepo.Elem().Interface().(T)
	} else {
		err := fmt.Errorf("*TenantDb field not found in type %s", typ.String())

		ret := createErrorRepoFromType(typ, err)
		return ret.(T)
	}
}

// func RepositoryFromType(db *sql.DB, typ reflect.Type) reflect.Value {
// 	tableName:=internal.Utils.CacheTableNameFromStruct
// 	ret := reflect.New(typ).Elem()
// 	for i := 0; i < typ.NumField(); i++ {
// 		field := typ.Field(i)
// 		internal.EntityUtils.QueryableFromType(db, field.Type, ret.Field(i))

// 	}
// }
