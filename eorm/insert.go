package eorm

import (
	"database/sql"
	"eorm/migrate"
	"eorm/tenantDB"
	"fmt"
	"reflect"
)

type entityInfo struct {
	tableName string
	entity    *migrate.Entity
}
type inserter struct {
}

func (r inserter) getEntityInfo(typ reflect.Type) *entityInfo {
	model := ModelRegistry.GetModelByType(typ)
	tableName := model.GetTableName()
	entity := model.GetEntity()
	return &entityInfo{
		tableName: tableName,
		entity:    &entity,
	}
}

func (r inserter) fetchAfterInsert(dialect Dialect, res sql.Result, entity *migrate.Entity, dataValue reflect.Value) error {
	// Nếu không có cột tự tăng thì bỏ qua
	autoCols := entity.GetAutoValueColumns()
	if len(autoCols) == 0 {
		return nil
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return dialect.ParseError(err)
	}

	for _, col := range autoCols {
		valField := dataValue.FieldByName(col.Field.Name)
		if valField.CanSet() {
			switch valField.Kind() {
			case reflect.Int, reflect.Int64, reflect.Uint, reflect.Uint64:
				valField.SetInt(lastID)
			default:
				return fmt.Errorf("unsupported auto-increment type: %s", valField.Kind())
			}
		}
	}

	return nil
}

func (r *inserter) InsertWithTx(tx *tenantDB.TenantTx, data interface{}) error {
	dialect := dialectFactory.Create(tx.Db.GetDriverName())
	dataValue := reflect.ValueOf(data)
	typ := reflect.TypeOf(data)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		dataValue = dataValue.Elem()
	}
	repoType := r.getEntityInfo(typ)
	sql, args := dialect.MakeSqlInsert(repoType.tableName, repoType.entity.GetColumns(), data)
	sqlStmt, err := tx.Prepare(sql)
	if err != nil {
		return err
	}
	defer sqlStmt.Close()
	sqlResult, err := sqlStmt.Exec(args...)
	if err != nil {
		errParse := dialect.ParseError(err)
		if errParse, ok := errParse.(*DialectError); ok {
			errParse.Tables = []string{repoType.tableName}
			errParse.Fields = []string{repoType.entity.GetFieldByColumnName(errParse.DbCols[0])}
		}
		return errParse
	}

	err = r.fetchAfterInsert(dialect, sqlResult, repoType.entity, dataValue)
	if err != nil {
		if dialectError, ok := err.(*DialectError); ok {
			dialectError.Tables = []string{repoType.tableName}
			dialectError.Fields = []string{repoType.entity.GetFieldByColumnName(dialectError.DbCols[0])}
		}
		return err
	}
	return nil
}
func (r *inserter) Insert(db *tenantDB.TenantDB, data interface{}) error {
	dialect := dialectFactory.Create(db.GetDriverName())
	dataValue := reflect.ValueOf(data)
	typ := reflect.TypeOf(data)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		dataValue = dataValue.Elem()
	}
	repoType := r.getEntityInfo(typ)
	sql, args := dialect.MakeSqlInsert(repoType.tableName, repoType.entity.GetColumns(), data)
	sqlStmt, err := db.Prepare(sql)
	if err != nil {
		return err
	}
	defer sqlStmt.Close()
	sqlResult, err := sqlStmt.Exec(args...)
	if err != nil {
		errParse := dialect.ParseError(err)
		if errParse, ok := errParse.(*DialectError); ok {
			errParse.Tables = []string{repoType.tableName}
			errParse.Fields = []string{repoType.entity.GetFieldByColumnName(errParse.DbCols[0])}
		}
		return errParse
	}

	err = r.fetchAfterInsert(dialect, sqlResult, repoType.entity, dataValue)
	if err != nil {
		if dialectError, ok := err.(*DialectError); ok {
			dialectError.Tables = []string{repoType.tableName}
			dialectError.Fields = []string{repoType.entity.GetFieldByColumnName(dialectError.DbCols[0])}
		}
		return err
	}
	return nil
}

var inserterObj = &inserter{}

func assertSinglePointerToStruct(obj interface{}) error {
	v := reflect.ValueOf(obj)
	t := v.Type()

	depth := 0
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
		depth++
	}

	if depth != 1 {
		return fmt.Errorf("Insert: expected pointer to struct (*T), got %d-level pointer", depth)
	}

	if t.Kind() != reflect.Struct {
		return fmt.Errorf("Insert: expected pointer to struct, got pointer to %s", t.Kind())
	}

	return nil
}
func Insert(db *tenantDB.TenantDB, data interface{}) error {
	err := assertSinglePointerToStruct(data)
	if err != nil {
		return err
	}
	m, err := NewMigrator(db)
	if err != nil {
		return err
	}
	err = m.DoMigrates()
	if err != nil {
		return err
	}

	return inserterObj.Insert(db, data)
}
func InsertWithTx(tx *tenantDB.TenantTx, data interface{}) error {
	err := assertSinglePointerToStruct(data)
	if err != nil {
		return err
	}
	m, err := NewMigrator(tx.Db)
	if err != nil {
		return err
	}
	err = m.DoMigrates()
	if err != nil {
		return err
	}

	return inserterObj.InsertWithTx(tx, data)
}
func InsertBatch[T any](db *tenantDB.TenantDB, data []T) (int64, error) {
	m, err := NewMigrator(db)
	if err != nil {
		return 0, err
	}
	err = m.DoMigrates()
	if err != nil {
		return 0, err
	}

	dialect := dialectFactory.Create(db.GetDriverName())
	repoType := inserterObj.getEntityInfo(reflect.TypeOf(data[0]))
	sql, args := dialect.MakeSqlInsertBatch(repoType.tableName, repoType.entity.GetColumns(), data)
	sqlStmt, err := db.Prepare(sql)
	if err != nil {
		errParse := dialect.ParseError(err)
		if errParse, ok := errParse.(*DialectError); ok {
			errParse.Tables = []string{repoType.tableName}
			errParse.Fields = []string{repoType.entity.GetFieldByColumnName(errParse.DbCols[0])}
		}
		return 0, errParse
	}
	defer sqlStmt.Close()
	sqlResult, err := sqlStmt.Exec(args...)
	if err != nil {
		errParse := dialect.ParseError(err)
		if errParse, ok := errParse.(*DialectError); ok {
			errParse.Tables = []string{repoType.tableName}
			errParse.Fields = []string{repoType.entity.GetFieldByColumnName(errParse.DbCols[0])}
		}
		return 0, errParse
	}
	rowsAff, err := sqlResult.RowsAffected()
	if err != nil {
		errParse := dialect.ParseError(err)
		if errParse, ok := errParse.(*DialectError); ok {
			errParse.Tables = []string{repoType.tableName}
			errParse.Fields = []string{repoType.entity.GetFieldByColumnName(errParse.DbCols[0])}
		}
		return 0, errParse
	}

	return rowsAff, nil
}
