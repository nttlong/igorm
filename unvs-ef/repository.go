package unvsef

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

func buildRepositoryFromStruct[T any](db *sql.DB, autoMigrate bool) (*T, error) {
	var v T
	typ := reflect.TypeOf(v)
	if typ == nil {
		typ = reflect.TypeOf((*T)(nil)).Elem()
	}
	tenantDb, err := utils.getTenantDb(db, typ)
	if err != nil {
		return nil, err
	}
	if tenantDb == nil {
		example := `type YourSchema struct {
						TenantDb
						}`
		return nil, fmt.Errorf("Repository struct '%s' must have embedded TenantDb struct, looks like this:\n:%s", typ.String(), example)
	}

	ret, err := utils.buildRepositoryFromType(typ)
	if err != nil {
		return nil, err
	}

	if len(ret.EntityTypes) == 0 {
		example := `\n  type User struct {
						Id   DbField[uint64] 'db:"primaryKey;autoIncrement"'
						Code DbField[string] 'db:"length(50)"'
					}`
		example = strings.ReplaceAll(example, "''", "`")

		return nil, fmt.Errorf("no entity type found in %s,'%s' must have at least one entity type looks like this\n:%s", typ.String(), typ.String(), example)
	}
	sqlMigrates, err := utils.GetScriptMigrate(db, tenantDb.DbName, tenantDb.Dialect, ret.EntityTypes...)
	if err != nil {
		return nil, err
	}
	tenantDb.SqlMigrate = sqlMigrates
	if autoMigrate {
		err = tenantDb.DoMigrate()
		if err != nil {
			return nil, err
		}
	}
	ret.ValueOfRepo.FieldByName("TenantDb").Set(reflect.ValueOf(*tenantDb))

	retVal := ret.ValueOfRepo.Interface().(T)

	return &retVal, nil

}

var repoCache sync.Map

func Repo[T any](db *sql.DB, autoMigrate bool) (*T, error) {
	typ := reflect.TypeOf((*T)(nil)).Elem()
	dbName, err := utils.GetDbName(db)
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("%s_%s", dbName, typ.String())
	if v, ok := repoCache.Load(key); ok {
		return v.(*T), nil
	}

	ret, err := buildRepositoryFromStruct[T](db, autoMigrate)
	if err != nil {
		return nil, err
	}

	repoCache.Store(key, ret)
	return ret, nil
}
