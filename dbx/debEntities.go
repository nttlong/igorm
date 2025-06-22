package dbx

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

func getStructFieldValue(s interface{}, fieldName string) (interface{}, error) {
	val := reflect.ValueOf(s)

	// Ensure it's a struct or a pointer to a struct
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input is not a struct or a pointer to a struct")
	}

	field := val.FieldByName(fieldName)
	if !field.IsValid() {
		return nil, fmt.Errorf("field '%s' not found in struct", fieldName)
	}

	return field.Interface(), nil
}
func createInsertCommand2(entity interface{}, entityType *EntityType) (*sqlWithParams, error) {
	var ret = sqlWithParams{
		Params: []interface{}{},
	}

	ret.Sql = "insert into "
	fields := []string{}
	valParams := []string{}
	// fields := getAllFields(typ)
	for _, field := range entityType.EntityFields {

		if field.IsPrimaryKey && field.DefaultValue == "auto" {
			continue

		}

		fieldVal, err := getStructFieldValue(entity, field.Name)
		if err != nil {
			return nil, err
		}
		if fieldVal == nil && !field.AllowNull && field.DefaultValue == "" {
			if val, ok := mapDefaultValueOfGoType[field.NonPtrFieldType]; ok {
				ret.Params = append(ret.Params, val)
				fields = append(fields, field.Name)
				valParams = append(valParams, "?")
			}
		} else {
			ret.Params = append(ret.Params, fieldVal)
			fields = append(fields, field.Name)
			valParams = append(valParams, "?")
		}

	}
	ret.Sql += entityType.TableName + " (" + strings.Join(fields, ",") + ") values (" + strings.Join(valParams, ",") + ")"
	return &ret, nil
}

var getSlqInsertCache = sync.Map{}

func getSlqInsert(entityType *EntityType) string {
	if v, ok := getSlqInsertCache.Load(entityType.TableName); ok {
		return v.(string)
	}
	retSql := "insert into "
	fields := []string{}
	valParams := []string{}
	for _, field := range entityType.EntityFields {

		if field.IsPrimaryKey && field.DefaultValue == "auto" {
			continue

		}
		fields = append(fields, field.Name)
		valParams = append(valParams, "?")

	}
	retSql += entityType.TableName + " (" + strings.Join(fields, ",") + ") values (" + strings.Join(valParams, ",") + ")"
	getSlqInsertCache.Store(entityType.TableName, retSql)
	return retSql
}
func createInsertCommand(entity interface{}, entityType *EntityType) (*sqlWithParams, error) {
	var ret = sqlWithParams{
		Params: []interface{}{},
	}

	//fields := []string{}
	//valParams := []string{}
	// // fields := getAllFields(typ)
	// for _, field := range entityType.EntityFields {

	// 	if field.IsPrimaryKey && field.DefaultValue == "auto" {
	// 		continue

	// 	}
	// 	fields = append(fields, field.Name)
	// 	valParams = append(valParams, "?")

	// }
	//ret.Sql = entityType.TableName + " (" + strings.Join(fields, ",") + ") values (" + strings.Join(valParams, ",") + ")"
	// fields := getAllFields(typ)
	ret.Sql = getSlqInsert(entityType)
	for _, field := range entityType.EntityFields {

		if field.IsPrimaryKey && field.DefaultValue == "auto" {
			continue

		}

		fieldVal, err := getStructFieldValue(entity, field.Name)
		if err != nil {
			return nil, err
		}
		if fieldVal == nil && !field.AllowNull && field.DefaultValue == "" {
			if val, ok := mapDefaultValueOfGoType[field.NonPtrFieldType]; ok {
				ret.Params = append(ret.Params, val)

			}
		} else {
			ret.Params = append(ret.Params, fieldVal)
		}

	}
	// ret.Sql += entityType.TableName + " (" + strings.Join(fields, ",") + ") values (" + strings.Join(valParams, ",") + ")"
	return &ret, nil
}
func Count[T any](ctx *DBXTenant, where string, args ...interface{}) (int64, error) {
	entityType := reflect.TypeFor[T]()
	e, err := Entities.CreateEntityType(entityType)
	if err != nil {
		return 0, err
	}
	sqlCount := "select count(*) as count from " + e.TableName
	if where != "" {
		sqlCount += " where " + where
	}
	execSQl, err := ctx.compiler.Parse(sqlCount, args...)
	if err != nil {
		return 0, err
	}
	var count int64
	err = ctx.DB.QueryRow(execSQl, args...).Scan(&count)
	if err != nil {
		if dbxErr := MssqlErrorParser.ParseError(nil, ctx.DB, err); dbxErr != nil {
			return 0, dbxErr
		}
		return 0, err
	}
	return count, nil

}
func CountWithContext[T any](ctx context.Context, db *DBXTenant, where string, args ...interface{}) (int64, error) {
	entityType := reflect.TypeFor[T]()
	e, err := Entities.CreateEntityType(entityType)
	if err != nil {
		return 0, err
	}
	sqlCount := "select count(*) as count from " + e.TableName
	if where != "" {
		sqlCount += " where " + where
	}
	var count int64
	execSQl, err := db.compiler.Parse(sqlCount, args...)
	if err != nil {
		return 0, err
	}
	err = db.DB.QueryRowContext(ctx, execSQl, args...).Scan(&count)
	if err != nil {
		if dbxErr := MssqlErrorParser.ParseError(ctx, db.DB, err); dbxErr != nil {
			return 0, dbxErr
		}
		return 0, err
	}
	return count, nil

}
