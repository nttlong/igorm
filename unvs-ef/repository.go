package unvsef

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

/*
While create repository looks like

	type Repository struct {
		*TenantDb //<-- Should I user pointer here
		Articles  *Article
		Comments  *Comment
	}
	func (r *Repository) Init() { // the relationship set up herr
			r.NewRelationship().From(r.Articles.Id).To(r.Comments.ArticleId, r.Comments.Id)

	}
*/
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
						TenantDb //or *TenantDb
						}`
		return nil, fmt.Errorf("Repository struct '%s' must have embedded TenantDb struct, looks like this:\n:%s", typ.String(), example)
	}

	ret, err := utils.buildRepositoryFromType(typ)
	if err != nil {
		if repositoryErr, ok := err.(buildRepositoryError); ok {
			example := `type ` + repositoryErr.FieldName + ` struct {
					Entity[` + repositoryErr.FieldName + `] // or Entity[` + repositoryErr.FieldName + `] 'db:"table(orders)"' if you want to specify table name
				}`
			example = strings.ReplaceAll(example, "'", "`")
			return nil, fmt.Errorf("build repository error for %s: embedded Entity was not found in    %s \n\t\tclarification looks like this:\n:%s", typ.String(), repositoryErr.FieldTypeName, example)
		}
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
	if autoMigrate {
		sqlMigrates, err := utils.GetScriptMigrate(db, tenantDb.DbName, tenantDb.Dialect, ret.EntityTypes...)
		if err != nil {
			return nil, err
		}
		tenantDb.SqlMigrate = sqlMigrates

		err = tenantDb.DoMigrate()
		if err != nil {
			return nil, err
		}
	}
	tenantDbField := ret.ValueOfRepo.FieldByName("TenantDb")
	if tenantDbField.Kind() == reflect.Ptr { //<-- important!: if repository struct has pointer to TenantDb, we need to set it to point to the actual TenantDb struct
		// tenantDbField = tenantDbField.Elem()
		tenantDbField.Set(reflect.ValueOf(tenantDb))
	} else { //<-- if repository struct has TenantDb struct directly, we need to set it to point to the actual TenantDb struct
		tenantDbField.Set(reflect.ValueOf(*tenantDb))
	}

	method := ret.PtrValueOfRepo.MethodByName("Init")
	if method.IsValid() && method.Type().NumIn() == 0 { //<-- check has init function for Repository
		method.Call([]reflect.Value{})
	}
	if tenantDb.buildRepositoryError != nil {
		return nil, tenantDb.buildRepositoryError
	}
	retVal := ret.ValueOfRepo.Interface().(T)

	return &retVal, nil

}

var repoCache sync.Map

/*
While create repository looks like

	type Repository struct {
		*TenantDb //<-- Should I user pointer here
		Articles  *Article
		Comments  *Comment
	}
		func (r *Repository) Init() { // the relationship set up herr
				r.NewRelationship().From(r.Articles.Id).To(r.Comments.ArticleId, r.Comments.Id)

		}
*/
func createErrorRepo(typ reflect.Type, err error) interface{} {
	if tenantFieldStruct, ok := typ.FieldByName("TenantDb"); ok {

		typeOfTenantFieldStruct := tenantFieldStruct.Type
		if typeOfTenantFieldStruct.Kind() == reflect.Ptr {
			typeOfTenantFieldStruct = typeOfTenantFieldStruct.Elem()
		}

		ret := reflect.New(typ).Elem()

		tenantField := ret.FieldByName("TenantDb")
		if tenantField.Kind() == reflect.Ptr {
			tenantFieldVal := reflect.New(typeOfTenantFieldStruct)
			tenantFieldIns := tenantFieldVal.Elem()
			tenantFieldIns.FieldByName("Err").Set(reflect.ValueOf(err))
			tenantField.Set(tenantFieldIns.Addr())
		} else {
			tenantField.FieldByName("Err").Set(reflect.ValueOf(err))

		}

		return ret.Interface()
	} else {
		panic(fmt.Sprintf("Repository struct '%s' must have embedded TenantDb struct", typ.String()))
	}
}
func Repo[T any](db *sql.DB, autoMigrate bool) *T {
	typ := reflect.TypeOf((*T)(nil)).Elem()
	if db == nil {
		err := fmt.Errorf("db is nil")
		ret := createErrorRepo(typ, err).(T)
		return &ret

	}
	dbName, err := utils.GetDbName(db)

	if err != nil { //<-- error getting db name return  repository with error status

		ret := createErrorRepo(typ, err).(T)
		return &ret

	}
	key := fmt.Sprintf("%s_%s", dbName, typ.String())
	if v, ok := repoCache.Load(key); ok {
		return v.(*T)
	}

	ret, err := buildRepositoryFromStruct[T](db, autoMigrate)
	if err != nil {
		ret := createErrorRepo(typ, err).(T)
		return &ret
	}

	repoCache.Store(key, ret)
	return ret
}
