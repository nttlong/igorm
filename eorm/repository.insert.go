package eorm

import (
	"context"
	"database/sql"
	"eorm/migrate"
	"eorm/tenantDB"
	"reflect"
)

type Repository[T any] struct {
	tableName string
	entity    *migrate.Entity
}

func (r *Repository[T]) fetchAfterInsert(dialect Dialect, sqlRow *sql.Row, entity *T) (*T, error) {
	if sqlRow.Err() != nil {
		return entity, sqlRow.Err()
	}
	scanReData := []interface{}{}
	valOfEntity := reflect.ValueOf(entity).Elem()
	for _, col := range r.entity.GetAutoValueColumns() {
		valField := valOfEntity.FieldByName(col.Field.Name)
		if valField.CanAddr() && valField.CanSet() {
			scanReData = append(scanReData, valField.Addr().Interface())
		}
	}
	err := sqlRow.Scan(scanReData...)
	if err != nil {

		return entity, dialect.ParseError(err)
	}

	return entity, nil
}
func (r *Repository[T]) InsertContext(context context.Context, db *tenantDB.TenantDB, entity *T) (*T, error) {
	dialect := dialectFactory.Create(db.GetDriverName())

	sql, args := dialect.MakeSqlInsert(r.tableName, r.entity.GetColumns(), entity)
	sqlStmt, err := db.Prepare(sql)
	if err != nil {
		return entity, err
	}
	defer sqlStmt.Close()
	sqlRow := sqlStmt.QueryRowContext(context, args...)

	return r.fetchAfterInsert(dialect, sqlRow, entity)

}
func (r *Repository[T]) Insert(db *tenantDB.TenantDB, entity *T) (*T, error) {
	dialect := dialectFactory.Create(db.GetDriverName())

	sql, args := dialect.MakeSqlInsert(r.tableName, r.entity.GetColumns(), entity)
	sqlStmt, err := db.Prepare(sql)
	if err != nil {
		return entity, err
	}
	defer sqlStmt.Close()
	sqlRow := sqlStmt.QueryRow(args...)

	return r.fetchAfterInsert(dialect, sqlRow, entity)
}
func (r *Repository[T]) InsertWithTx(tx *tenantDB.TenantTx, entity *T) (*T, error) {

	dialect := dialectFactory.Create(tx.GetDriverName())
	sql, args := dialect.MakeSqlInsert(r.tableName, r.entity.GetColumns(), entity)
	sqlStmt, err := tx.Prepare(sql)
	if err != nil {
		return entity, err
	}
	defer sqlStmt.Close()
	sqlRow := sqlStmt.QueryRow(args...)

	return r.fetchAfterInsert(dialect, sqlRow, entity)

}

type OnQuoteFunc = func(string, string) string

var OnQuote OnQuoteFunc
