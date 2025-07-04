package orm

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	internal "unvs-orm/internal"
)

type fx struct {
	colName string
	field   internal.FieldTag
}
type Object[T any] struct {
	tenantDb *internal.TenantDb
	Data     T
	stmt     *sql.Stmt
	meta     []fx
}

func (e *Object[T]) Use(tx *sql.Tx) error {
	tableName := internal.Utils.TableNameFromStruct(reflect.TypeFor[T]())
	metaData := internal.Utils.GetMetaInfo(reflect.TypeFor[T]())
	e.meta = []fx{}
	for _, meta := range metaData {

		args := []interface{}{}
		dataVal := reflect.ValueOf(e.Data)

		for k, v := range meta {

			if v.AutoIncrement {
				fmt.Println("auto increment field found, skipping")
				continue
			}
			field := dataVal.FieldByName(v.Field.Name)

			if !field.IsValid() || field.IsZero() {
				continue
			}
			valField := field.FieldByName("Val")
			if !valField.IsValid() || valField.IsZero() {
				args = append(args, nil)
			} else {
				args = append(args, valField.Interface())
			}
			e.meta = append(e.meta, fx{
				colName: k,
				field:   v,
			})

		}

	}
	sql := "INSERT INTO " + tableName
	fields := []string{}
	values := []string{}

	for _, v := range e.meta {

		fields = append(fields, v.colName)
		values = append(values, "?")
	}

	sql = sql + "(" + strings.Join(fields, ",") + ") VALUES (" + strings.Join(values, ",") + ")"

	stmt, err := tx.Prepare(sql)
	if err != nil {
		return err
	}
	e.stmt = stmt
	return nil
}
func (e *Object[T]) InsertWithTransaction(tx *sql.Tx) error {

	args := []interface{}{}
	dataVal := reflect.ValueOf(e.Data)

	for _, v := range e.meta {

		field := dataVal.FieldByName(v.field.Field.Name)

		if !field.IsValid() || field.IsZero() {
			fmt.Println("auto increment field found, skipping")
			continue
		}
		valField := field.FieldByName("Val")
		if !valField.IsValid() || valField.IsZero() {
			args = append(args, nil)
		} else {
			args = append(args, valField.Interface())
		}

	}
	//"INSERT INTO orders(updated_by,version,note,created_at,updated_at,created_by) VALUES (@UpdatedBy,@Version,@Note,@CreatedAt,@UpdatedAt,@CreatedBy)"

	_, err := e.stmt.Exec(args...)
	if err != nil {
		return err
	} else {
		return nil
	}

}
func (e *Object[T]) Insert() error {
	tableName := internal.Utils.TableNameFromStruct(reflect.TypeFor[T]())
	meta := internal.Utils.GetMetaInfo(reflect.TypeFor[T]())

	for _, meta := range meta {
		sql := "INSERT INTO " + tableName
		fields := []string{}
		values := []string{}
		args := []interface{}{}
		dataVal := reflect.ValueOf(e.Data)

		for k, v := range meta {

			if v.AutoIncrement {
				continue
			}
			field := dataVal.FieldByName(v.Field.Name)

			if !field.IsValid() || field.IsZero() {
				continue
			}
			valField := field.FieldByName("Val")
			if !valField.IsValid() || valField.IsZero() {
				args = append(args, nil)
			} else {
				args = append(args, valField.Interface())
			}

			fields = append(fields, k)
			values = append(values, "?")
		}

		sql = sql + "(" + strings.Join(fields, ",") + ") VALUES (" + strings.Join(values, ",") + ")"

		_, err := e.tenantDb.Exec(sql, args...)
		return err
	}
	return fmt.Errorf("no table found for entity %s", reflect.TypeFor[T]().String())
}
