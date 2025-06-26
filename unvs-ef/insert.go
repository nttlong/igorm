package unvsef

import (
	"fmt"
	"reflect"
	"strings"
)

type InsertQuery struct {
	Entity any
}

func Insert(entity any) *InsertQuery {
	return &InsertQuery{Entity: entity}
}
func (q *InsertQuery) ToSQL(d Dialect) (string, []interface{}) {
	val := reflect.ValueOf(q.Entity)
	typ := reflect.TypeOf(q.Entity)

	if val.Kind() == reflect.Pointer {
		val = val.Elem()
		typ = typ.Elem()
	}

	tableName := utils.TableNameFromStruct(typ)
	meta := utils.GetMetaInfo(typ)

	columns := []string{}
	placeholders := []string{}
	args := []interface{}{}

	for fieldName, tag := range meta[tableName] {
		field := val.FieldByName(tag.Field.Name)
		if field.Kind() == reflect.Struct {
			// field is DbField[T], check Value field
			v := field.FieldByName("Value")
			if !v.IsNil() {
				columns = append(columns, d.QuoteIdent(tableName, fieldName))
				placeholders = append(placeholders, "?")
				args = append(args, v.Elem().Interface())
			}
		}
	}

	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		d.QuoteIdent(tableName),
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	return sql, args
}
