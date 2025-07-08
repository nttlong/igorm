package internal

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"unicode"

	pluralizeLib "github.com/gertd/go-pluralize"
)

var pluralize = pluralizeLib.NewClient()

type utilsPackage struct {
	cacheGetMetaInfo                        sync.Map
	CacheTableNameFromStruct                sync.Map
	cacheGetPkFromMeta                      sync.Map
	cacheGetUniqueConstraintsFromMetaByType sync.Map
	cacheGetIndexConstraintsFromMetaByType  sync.Map
	schemaCache                             sync.Map
	cacheGetOrCreateRepository              sync.Map
	cacheGetTenantDb                        sync.Map
	cacheBuildFieldMap                      sync.Map
	mapType                                 map[reflect.Type]string
	currentPackagePath                      string //<-- cache current package path
	cacheGetRequireFields                   sync.Map
	cacheGetAutoPkKey                       sync.Map
	entityTypeName                          string
	cacheVerifyModelFieldFirst              sync.Map

	cacheReplacePlaceHolder          sync.Map
	cacheEntityNameAndDbTableName    sync.Map
	cacheEntityNameAndMetaInfo       sync.Map
	cacheToSnakeCase                 sync.Map
	cacheGetTableNameFromVirtualName sync.Map
}

/*
	   Get Go type of DbField in name

	   Example:
	    var testType= DbField[time.Time]{}
		utils.ResolveFieldKind(reflect.TypeOf(&testType))-> "time.Time"
*/
func (u *utilsPackage) ResolveFieldKind(field reflect.StructField) string {

	ftType := field.Type
	if ftType.Kind() == reflect.Ptr {
		ftType = ftType.Elem()
	}

	if ftType.PkgPath() == u.currentPackagePath {

		if ftType.Kind() == reflect.Ptr {
			ftType = ftType.Elem()
		}
		if typeName, ok := u.mapType[ftType]; ok {
			return typeName
		} else {
			panic(fmt.Errorf("'utilsPackage.ResolveFieldKind (row 204)' report: %s was not found in mapType of utilsPackage", ftType.String()))
		}
	}
	strFt := field.Type.String()
	if strings.Contains(strFt, ".DbField[") {
		typeParam := strings.Split(strFt, ".DbField[")[1]
		typeParam = strings.Split(typeParam, "]")[0]
		typeParam = strings.TrimLeft(typeParam, "*")
		return typeParam
	}
	return field.Type.Kind().String()
}

/*
The purpose support for Dialects is to provide a way to customize the SQL generation for different databases.
*/
func (u *utilsPackage) Quote(strQuote string, str ...string) string {
	left := strQuote[0:1]
	right := strQuote[1:2]
	ret := left + strings.Join(str, left+"."+right) + right
	return ret
}

/*
The function will get all primary key cols from type (the  info will be stored in cache . The next call will return form cache)
return

	{
		"<primary key constraint name": {//<-- The value is obtained by the combination of the table name, double underscores ("__"), and the column name.
			"<key field name 1>": {...},//<--FieldTag Info
			...
			"<key field name n>": {...},//<--FieldTag Info
		}
	}
*/
func (u *utilsPackage) GetPkFromMetaByType(typ reflect.Type) map[string]map[string]FieldTag {
	// 1. Kiểm tra cache trước (check cache first)
	if pk, ok := u.cacheGetPkFromMeta.Load(typ); ok {
		return pk.(map[string]map[string]FieldTag)
	}
	metaInfo := u.GetMetaInfo(typ)
	ret := make(map[string]map[string]FieldTag)
	fieldsNames := []string{}
	for tableName, fields := range metaInfo {
		pkMap := make(map[string]FieldTag)
		for fieldName, fieldTag := range fields {
			if fieldTag.PrimaryKey {
				pkMap[fieldName] = fieldTag
				fieldsNames = append(fieldsNames, fieldName)

			}
		}
		ret[tableName+"_"+strings.Join(fieldsNames, "_")] = pkMap

	}
	u.cacheGetPkFromMeta.Store(typ, ret)
	return ret
}

func (u *utilsPackage) ToSnakeCase(str string) string {
	if v, ok := u.cacheToSnakeCase.Load(str); ok {
		return v.(string)
	}
	var result []rune
	for i, r := range str {
		if i > 0 && unicode.IsUpper(r) &&
			(unicode.IsLower(rune(str[i-1])) || (i+1 < len(str) && unicode.IsLower(rune(str[i+1])))) {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}
	ret := string(result)
	u.cacheToSnakeCase.Store(str, ret)
	return ret
}

func (u *utilsPackage) extractName(s string) string {
	re := regexp.MustCompile(`\((.*?)\)`)
	matches := re.FindStringSubmatch(s)
	if len(matches) == 2 {
		return matches[1]
	}
	return ""
}

/*
Get table name from struct

	 Example :

	 type User struct {
		  _ DbField[any] `db:table(my_user)`

	 }

	 return my_user

	 type User struct {} -> "users"
*/
func (u *utilsPackage) TableNameFromStruct(typ reflect.Type) string {
	// Check for override via table(...) tag
	if v, ok := u.CacheTableNameFromStruct.Load(typ); ok {
		return v.(string)
	}
	for i := 0; i < typ.NumField(); i++ {
		fType := typ.Field(i).Type
		if fType.Kind() == reflect.Ptr {
			fType = fType.Elem()
		}

		if strings.HasPrefix(fType.String(), u.entityTypeName) {
			parsed := u.ParseDBTag(typ.Field(i))
			if parsed.TableName != "" {
				u.CacheTableNameFromStruct.Store(typ, parsed.TableName)
				u.cacheEntityNameAndDbTableName.Store(strings.ToLower(typ.Name()), parsed.TableName)
				u.cacheEntityNameAndDbTableName.Store(u.Plural(strings.ToLower(typ.Name())), parsed.TableName)
				return parsed.TableName
			} else {
				ret := pluralize.Plural(u.ToSnakeCase(typ.Name()))
				u.CacheTableNameFromStruct.Store(typ, ret)
				u.cacheEntityNameAndDbTableName.Store(strings.ToLower(typ.Name()), ret)
				u.cacheEntityNameAndDbTableName.Store(u.Plural(strings.ToLower(typ.Name())), ret)
				return ret

			}
		}
	}
	return ""
}
func (u *utilsPackage) Plural(txt string) string {
	return pluralize.Plural(txt)
}

/*
The method will get all Unique Constraint from declarification type
#Exmaple

	type struct UserRole {
			Id   DbField[uint64]     `db:"primaryKey;autoIncrement"`
			UserId  DbField[uint64]  `db:"unique(user_role_idx)"` //<-- unique constraint 2 columns
			RoleId  DbField[uint64]  `db:"unique(user_role_idx)"` //<-- unique constraint 2 columns
			OwnerId DbField[uint64] `db:"unique"` //<-- unique constraint 1 column
	}

	return {
		"user_roles":{ //<-- table name
			"user_role_idx____user_roles___user_id__role_id":{ //<--- convention of constraint name is the combination of constraint name in tag, four underscores and table name (snake case)
				"user_id": {...} //<-- col 1
				"role_id": {...} // <--col 2
			},
			"user_role_idx____user_roles__owner_id": { //<--- convention of constraint name is the combination of field name (snake case) four underscore and table name (snake case)
				"owner_id": {...} //<-- field name
			}

		}
	}
*/
func (u *utilsPackage) GetUniqueConstraintsFromMetaByType(typ reflect.Type) map[string]map[string]FieldTag {
	// 1. Kiểm tra cache trước (check cache first)
	if unique, ok := u.cacheGetUniqueConstraintsFromMetaByType.Load(typ); ok {
		return unique.(map[string]map[string]FieldTag)
	}
	metaInfo := u.GetMetaInfo(typ)
	info := make(map[string]map[string]FieldTag)
	tableName := u.TableNameFromStruct(typ)

	for _, fields := range metaInfo {
		for fieldName, fieldTag := range fields {

			if fieldTag.Unique {
				ukName := fieldTag.UniqueName //<-- use fieldTag.UniqueName can be field name if no unique index name in tag else name of unique index in tag
				if _, ok := info[ukName]; !ok {
					info[ukName] = make(map[string]FieldTag)

				}
				info[ukName][fieldName] = fieldTag

			}
		}
	}
	ret := make(map[string]map[string]FieldTag)
	for ukName, fields := range info {

		refFields := []string{}
		for fieldName := range fields {
			refFields = append(refFields, fieldName)
		}
		constraintName := ukName + "____" + tableName + "___" + strings.Join(refFields, "__")
		ret[constraintName] = fields
	}
	u.cacheGetUniqueConstraintsFromMetaByType.Store(typ, ret)
	return ret
}

/*
The method will get all Unique Constraint from declarification type
#Exmaple

	type struct UserRole {
			Id   DbField[uint64]     `db:"primaryKey;autoIncrement"`
			UserId  DbField[uint64]  `db:"index(user_role_idx)"` //<-- unique constraint 2 columns
			RoleId  DbField[uint64]  `db:"index(user_role_idx)"` //<-- unique constraint 2 columns
			OwnerId DbField[uint64] `db:"index"` //<-- unique constraint 1 column
	}

	return {
		"user_roles":{ //<-- table name
			"user_role_idx____user_roles":{ //<--- convention of constraint name is constraint name in tag +"___"+ table name
				"user_id": {...} //<-- col 1
				"role_id": {...} // <--col 2
			},
			"user_roles____owner_id_idx": { //<--- constraint name
				"owner_id": {...} //<-- field name
			}

		}
	}
*/
func (u *utilsPackage) GetIndexConstraintsFromMetaByType(typ reflect.Type) map[string]map[string]FieldTag {
	// 1. Kiểm tra cache trước (check cache first)
	if unique, ok := u.cacheGetIndexConstraintsFromMetaByType.Load(typ); ok {
		return unique.(map[string]map[string]FieldTag)
	}
	metaInfo := u.GetMetaInfo(typ)
	info := make(map[string]map[string]FieldTag)
	tableName := u.TableNameFromStruct(typ)

	for _, fields := range metaInfo {
		for fieldName, fieldTag := range fields {

			if fieldTag.Index {
				indexName := fieldTag.IndexName
				if _, ok := info[indexName]; !ok {
					info[indexName] = make(map[string]FieldTag)

				}
				info[indexName][fieldName] = fieldTag

			}
		}
	}
	ret := make(map[string]map[string]FieldTag)
	for idxName, fields := range info {

		refFields := []string{}
		for fieldName := range fields {
			refFields = append(refFields, fieldName)
		}
		constraintName := idxName + "____" + tableName + "___" + strings.Join(refFields, "__")
		ret[constraintName] = fields
	}
	u.cacheGetIndexConstraintsFromMetaByType.Store(typ, ret)
	return ret
}
func (u *utilsPackage) buildRepositoryFromType(typ reflect.Type, isChildren bool) (*repositoryValueStruct, error) {

	retValueOfRepo := reflect.New(typ)
	valueOfRepo := retValueOfRepo.Elem()

	baseType := reflect.TypeOf(Base{})
	entityTypes := []reflect.Type{}
	for i := 0; i < typ.NumField(); i++ {

		field := typ.Field(i)
		fieldType := field.Type
		if field.Anonymous {
			if fieldType.Kind() == reflect.Ptr {
				fieldType = fieldType.Elem()
			}

			if baseType.Name() == fieldType.Name() {
				if !isChildren {

					base := &Base{
						Relationships: []*RelationshipRegister{},
					}
					attachBase := OnRequestBaseFn(base)
					baseVal := reflect.ValueOf(attachBase)
					if valueOfRepo.Field(i).Kind() == reflect.Ptr {

						valueOfRepo.Field(i).Set(baseVal) //<-- "reflect: reflect.Value.Set using unaddressable value"
					} else {

						valueOfRepo.Field(i).Set(baseVal.Elem())
					}

				}
				continue

			} else {
				repoVal, err := u.buildRepositoryFromType(field.Type, true) //<-- do not gen sql migrate for inner entity
				if err != nil {
					return nil, err
				}
				entityTypes = append(entityTypes, repoVal.EntityTypes...)
				valueOfRepo.Field(i).Set(repoVal.ValueOfRepo.Addr())

				continue
			}
		}
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}
		tableName := utils.TableNameFromStruct(fieldType)
		if tableName == "" {

			return nil, buildRepositoryError{
				FieldName:     field.Name,
				FieldTypeName: fieldType.String(),
				err:           fmt.Errorf("embedded Model[] was not found struct %s of field %s of type %s is not a struct or not a pointer to struct", field.Name, fieldType),
			}
		}
		var modelVal *reflect.Value = nil
		queryableVal := entityUtils.QueryableFromType(fieldType, tableName, modelVal)

		queryableValField := valueOfRepo.Field(i)

		queryableValField.Set(queryableVal)

		entityType := field.Type
		if entityType.Kind() == reflect.Ptr {
			entityType = entityType.Elem()
		}
		entityTypes = append(entityTypes, entityType)
	}
	ret := &repositoryValueStruct{
		ValueOfRepo:    valueOfRepo,
		PtrValueOfRepo: retValueOfRepo,
		EntityTypes:    entityTypes,
	}
	return ret, nil
}
func (u *utilsPackage) verifyModelFieldFirst(typ reflect.Type) error {
	//check cache

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if _, ok := u.cacheVerifyModelFieldFirst.Load(typ); ok {
		return nil
	}

	if typ.NumField() == 0 || typ.Field(0).Name != "Model" {

		retErr := fmt.Errorf("orm.Model must be the first field in struct %s", typ.Name())
		u.cacheVerifyModelFieldFirst.Store(typ, retErr)
		return retErr
	}
	return nil
}
func (u *utilsPackage) GetOrCreateRepository(typ reflect.Type) (*repositoryValueStruct, error) {
	//check cache

	key := typ.String()
	if val, ok := u.cacheGetOrCreateRepository.Load(key); ok {
		return val.(*repositoryValueStruct), nil
	}
	repoVal, err := u.buildRepositoryFromType(typ, false)

	if err != nil {
		return nil, err
	}
	//set cache
	u.cacheGetOrCreateRepository.Store(key, repoVal)
	return repoVal, nil
}
func (u *utilsPackage) GetMetaInfoByTableName(tableName string) map[string]FieldTag {
	//check cache

	if val, ok := u.cacheEntityNameAndMetaInfo.Load(tableName); ok {
		return val.(map[string]FieldTag)
	}
	return nil
}

type FnGetBAse func() interface{}
type OnRequestBase func(base *Base) interface{}

var OnRequestBaseFn OnRequestBase
