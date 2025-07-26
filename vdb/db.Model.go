package vdb

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"
	"vdb/migrate"
	"vdb/tenantDB"
)

type dbModel struct {
	context     context.Context
	db          *TenantDB
	tableName   string
	cols        *[]migrate.ColumnDef
	mapCols     *map[string]migrate.ColumnDef
	where       string
	whereArgs   []interface{}
	fieldUpdate string
	valueUpdate interface{}
}
type UpdateResult struct {
	RowsAffected int64
	Error        error
	Sql          string //<-- if error is not nil, this field will be not empty
}
type ErrFieldNotFound struct {
	Field     string
	TableName string
}

func (e *ErrFieldNotFound) Error() string {
	return fmt.Sprintf("field %s not found in table %s", e.Field, e.TableName)
}

type initTenantDBModel struct {
	once sync.Once
	val  initTenantDBModelCacheItem
}
type initTenantDBModelCacheItem struct {
	tableName string
	cols      *[]migrate.ColumnDef
	keyCols   []migrate.ColumnDef
	mapCols   map[string]migrate.ColumnDef
	dialect   Dialect
}

var initTenantDBModelCache sync.Map

func (db *TenantDB) getModelFromCache(modelType reflect.Type) initTenantDBModelCacheItem {
	key := db.GetDriverName() + "://" + db.GetDBName() + "/" + modelType.String()
	actual, _ := initTenantDBModelCache.LoadOrStore(key, &initTenantDBModel{})
	init := actual.(*initTenantDBModel)
	init.once.Do(func() {
		repoType := inserterObj.getEntityInfo(modelType)
		tableName := repoType.tableName
		columns := repoType.entity.GetColumns()
		mapCols := make(map[string]migrate.ColumnDef)
		keyCols := []migrate.ColumnDef{}
		for _, col := range columns {
			mapCols[strings.ToLower(col.Name)] = col
			if col.PKName != "" {
				keyCols = append(keyCols, col)
			}
		}
		init.val = initTenantDBModelCacheItem{
			tableName: tableName,
			cols:      &columns,
			mapCols:   mapCols,
			dialect:   dialectFactory.create(db.GetDriverName()),
			keyCols:   keyCols,
		}
	})
	return init.val
}
func (db *TenantDB) Model(model interface{}) *dbModel {
	typ := reflect.TypeOf(model)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	modelInfo := db.getModelFromCache(typ)

	// dialect := dialectFactory.Create(db.GetDriverName())

	return &dbModel{
		db:        db,
		tableName: modelInfo.tableName,
		cols:      modelInfo.cols,
		mapCols:   &modelInfo.mapCols,
	}

}
func (m *dbModel) Context(context context.Context) *dbModel {
	m.context = context
	return m
}
func (m *dbModel) Where(expr string, args ...interface{}) *dbModel {
	m.where = expr
	m.whereArgs = args
	return m
}
func (m *dbModel) buildSqlNoCache() (string, error) {

	fieldMatched, ok := (*m.mapCols)[strings.ToLower(m.fieldUpdate)]
	if !ok {
		return "", fmt.Errorf("field %s not found in table %s", m.fieldUpdate, m.tableName)
	}
	dialect := dialectFactory.create(m.db.GetDriverName())
	source := dialect.Quote(m.tableName)
	if m.where != "" {
		compiler, err := NewExprCompiler(m.db.TenantDB)
		if err != nil {
			return "", err
		}

		compiler.context.purpose = build_purpose_where
		compiler.context.tables = []string{m.tableName}
		compiler.context.alias = map[string]string{m.tableName: m.tableName}
		compiler.context.paramIndex = 1
		err = compiler.buildWhere(m.where)
		if err != nil {
			return "", err
		}

		m.where = compiler.content

	}

	sql := ""
	strSet := fmt.Sprintf("%s = %s", dialect.Quote(fieldMatched.Name), dialect.ToParam(1))
	if m.where != "" {
		sql = fmt.Sprintf("UPDATE %s SET %s WHERE %s", source, strSet, m.where)
	} else {
		sql = fmt.Sprintf("UPDATE %s SET %s", source, strSet)
	}
	return sql, nil

}

type initBuildSql struct {
	once sync.Once
	val  string
	err  error
}

var buildSqlCache sync.Map

func (m *dbModel) BuildSql() (string, error) {
	key := m.db.GetDriverName() + "://" + m.db.GetDBName() + "/" + m.tableName + "/" + m.fieldUpdate + "/" + m.where
	actual, _ := buildSqlCache.LoadOrStore(key, &initBuildSql{})
	initBuild := actual.(*initBuildSql)
	initBuild.once.Do(func() {
		sql, err := m.buildSqlNoCache()
		initBuild.val = sql
		initBuild.err = err
	})
	return initBuild.val, initBuild.err

}

type initBuildSqlUpdateWithFieldsAndWhere struct {
	once sync.Once
	val  string
	err  error
}

var buildSqlUpdateWithFieldsAndWhereCache sync.Map

func buildSqlUpdateWithFieldsAndWhereWithCache(db *TenantDB, tableName string, fields []string, where string) (string, error) {
	key := db.GetDriverName() + "://" + db.GetDBName() + "/" + tableName + "/[" + strings.Join(fields, "][") + "]/" + where
	actual, _ := buildSqlUpdateWithFieldsAndWhereCache.LoadOrStore(key, &initBuildSqlUpdateWithFieldsAndWhere{})
	initBuild := actual.(*initBuildSqlUpdateWithFieldsAndWhere)
	initBuild.once.Do(func() {
		sql, err := buildSqlUpdateWithFieldsAndWhere(db, tableName, fields, where)
		initBuild.val = sql
		initBuild.err = err
	})
	return initBuild.val, initBuild.err
}

func buildSqlUpdateWithFieldsAndWhere(db *TenantDB, tableName string, fields []string, where string) (string, error) {
	dialect := dialectFactory.create(db.GetDriverName())
	source := dialect.Quote(tableName)
	compiler, err := NewExprCompiler(db.TenantDB)
	if err != nil {
		return "", err
	}

	compiler.context.tables = []string{tableName}
	compiler.context.alias = map[string]string{tableName: tableName}
	if where != "" {
		compiler.context.purpose = build_purpose_where
		compiler.context.paramIndex = 1
		err = compiler.buildWhere(where)
		if err != nil {
			return "", err
		}

		where = compiler.content

	}

	strSet := strings.Join(fields, ", ")
	compiler.context.purpose = build_purpose_for_update
	err = compiler.buildSetter(strSet)
	if err != nil {
		return "", err
	}

	strSet = compiler.content

	sql := ""
	if where != "" {
		sql = fmt.Sprintf("UPDATE %s SET %s WHERE %s", source, strSet, where)
	} else {
		sql = fmt.Sprintf("UPDATE %s SET %s", source, strSet)
	}
	return sql, nil

}
func (m *dbModel) updateByMap(data map[string]interface{}) (sql.Result, string, error) {
	args := []interface{}{}
	strFields := []string{}
	for k, v := range data {
		k = strings.TrimSuffix(strings.TrimPrefix(strings.ToLower(k), " "), " ")
		col, ok := (*m.mapCols)[k]
		if col.PKName != "" {
			continue
		}
		if !ok {
			return nil, "", fmt.Errorf("field %s not found in table %s", k, m.tableName)
		}
		dbField := ""
		if fn, ok := v.(dbFunCall); ok {
			args = append(args, fn.args...)
			dbField = col.Name + "=" + fn.expr

		} else {
			dbField = col.Name + "=?"
			args = append(args, v)
		}

		strFields = append(strFields, dbField)
	}

	sql, err := buildSqlUpdateWithFieldsAndWhereWithCache(m.db, m.tableName, strFields, m.where)

	if err != nil {
		return nil, sql, err
	}
	args = append(args, m.whereArgs...)
	start := time.Now()
	r, err := m.db.Exec(sql, args...)
	n := time.Since(start).Nanoseconds()
	fmt.Println("Update time:", n, "ns")
	if err != nil {
		return nil, sql, err
	}
	return r, sql, nil
	/*

	 */

}

type initBuldExrForUpdateFieldAndValue struct {
	once sync.Once
	val  string
	err  error
}

var buldExrForUpdateFieldAndValueCacheData sync.Map

func buldExrForUpdateFieldAndValueCache(db *tenantDB.TenantDB, tableName, expr string) (string, error) {
	key := db.GetDriverName() + "://" + db.GetDBName() + "/" + tableName + "/" + expr
	actual, _ := buldExrForUpdateFieldAndValueCacheData.LoadOrStore(key, &initBuldExrForUpdateFieldAndValue{})
	initBuild := actual.(*initBuldExrForUpdateFieldAndValue)
	initBuild.once.Do(func() {
		sql, err := buldExrForUpdateFieldAndValue(db, tableName, expr)
		initBuild.val = sql
		initBuild.err = err
	})
	return initBuild.val, initBuild.err
}
func buldExrForUpdateFieldAndValue(db *tenantDB.TenantDB, tableName, expr string) (string, error) {
	compiler, err := NewExprCompiler(db)
	if err != nil {
		return "", err
	}
	compiler.context.purpose = build_purpose_for_function
	compiler.context.tables = []string{tableName}
	compiler.context.alias = map[string]string{tableName: tableName}
	compiler.context.paramIndex = 1
	err = compiler.buildSelectField(expr)
	if err != nil {
		return "", err
	}
	return compiler.content, nil
}
func (m *dbModel) updateFieldAndValue(field string, value interface{}) (sql.Result, string, error) {
	m.fieldUpdate = field
	m.valueUpdate = value
	if fn, ok := value.(dbFunCall); ok {
		expr, err := buldExrForUpdateFieldAndValueCache(m.db.TenantDB, m.tableName, fn.expr)
		if err != nil {
			return nil, "", err
		}
		m.valueUpdate = expr
		sql, err := m.BuildSql()
		sql = strings.Replace(sql, "?", expr, 1)
		if err != nil {
			return nil, "", err
		}

		args := fn.args
		args = append(args, m.whereArgs...)
		if m.context != nil {
			ctx, cancel := context.WithTimeout(m.context, 5*time.Second)
			defer cancel()
			r, err := m.db.ExecContext(ctx, sql, args...)
			if err != nil {
				return nil, "", err
			}
			return r, sql, nil
		} else {
			r, err := m.db.Exec(sql, args...)
			if err != nil {
				return nil, "", err
			}
			return r, sql, nil
		}

	}
	sql, err := m.BuildSql()
	if err != nil {
		return nil, sql, err
	}

	args := []interface{}{value}
	args = append(args, m.whereArgs...)
	if m.context != nil {
		ctx, cancel := context.WithTimeout(m.context, 5*time.Second)
		defer cancel()
		r, err := m.db.ExecContext(ctx, sql, args...)
		return r, sql, err
	} else {
		r, err := m.db.Exec(sql, args...)
		return r, sql, err
	}

}
func (m *dbModel) parseUPdateError(result sql.Result, sql string, err error) UpdateResult {
	if err != nil {
		dialect := dialectFactory.create(m.db.GetDriverName())
		dError := dialect.ParseError(err)
		if dialectError, ok := dError.(*DialectError); ok {
			if dialectError.ConstraintName != "" {
				loader, errLoader := migrate.MigratorLoader(m.db.TenantDB)
				if errLoader == nil {
					schema, errLoader := loader.LoadFullSchema(m.db.TenantDB)
					if errLoader == nil {
						cols := schema.UniqueKeys[dialectError.ConstraintName]
						dialectError.Table = m.tableName

						for _, col := range cols.Columns {
							dialectError.DbCols = append(dialectError.Fields, col.Name)
							for _, col2 := range *m.cols {
								if col.Name == col2.Name {
									dialectError.Fields = append(dialectError.Fields, col2.Field.Name)
								}
							}

						}

						return UpdateResult{
							Error: dialectError,
							Sql:   sql,
						}
					}

				}

			}
		}

		return UpdateResult{
			Error: dError,
			Sql:   sql,
		}
	}
	r, err := result.RowsAffected()
	if err != nil {
		return UpdateResult{
			Error: err,
			Sql:   sql,
		}
	}
	if r == 0 {
		return UpdateResult{
			RowsAffected: r,
			Sql:          sql,
		}
	}
	return UpdateResult{
		RowsAffected: r,
	}
}
func (m *dbModel) Update(args ...interface{}) UpdateResult {
	if len(args) == 2 {
		if field, ok := args[0].(string); ok {
			r, sql, err := m.updateFieldAndValue(field, args[1])
			return m.parseUPdateError(r, sql, err)

		} else {
			panic("first argument must be a string")
		}
	} else if len(args) == 1 {
		if data, ok := args[0].(map[string]interface{}); ok {
			r, sql, err := m.updateByMap(data)
			return m.parseUPdateError(r, sql, err)
		} else {
			panic("first argument must be a map[string]interface{}")
		}
	}
	panic("invalid arguments, please use Update(field string, value interface{}) or Update(data map[string]interface{})")
}

type DeleteResult struct {
	RowsAffected int64
	Error        error
}

func (m *dbModel) Delete() DeleteResult {
	dialect := dialectFactory.create(m.db.GetDriverName())
	source := dialect.Quote(m.tableName)
	panic(source)
	if m.where != "" {
		compiler, err := NewExprCompiler(m.db.TenantDB)
		if err != nil {
			return DeleteResult{
				Error: err,
			}
		}

		compiler.context.purpose = build_purpose_where
		compiler.context.tables = []string{m.tableName}
		compiler.context.alias = map[string]string{m.tableName: m.tableName}
		compiler.context.paramIndex = 1
		err = compiler.buildWhere(m.where)
		if err != nil {
			return DeleteResult{
				Error: err,
			}
		}

		m.where = compiler.content
	}
	panic("not implemented")

}
