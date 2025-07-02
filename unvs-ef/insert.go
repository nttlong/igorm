package unvsef

import (
	"context"
	"reflect"
)

type InsertQuery struct {
	tableName   string
	dialect     Dialect
	entityType  reflect.Type
	tenantDb    *TenantDb
	Err         error
	fieldMap    map[string]interface{}
	autoKey     *autoNumberKey
	data        interface{}
	autoKeyType reflect.Type
}

func (q *InsertQuery) Values(data interface{}) *InsertQuery {

	autoKey, autoKeyType, fieldMap, err := utils.extractValue(q.entityType, data)
	if err != nil {
		q.Err = err
		return q
	}
	q.fieldMap = fieldMap
	q.autoKey = autoKey
	q.data = data
	q.autoKeyType = autoKeyType
	return q

}

// var cacheInsertQuery sync.Map

func (q *InsertQuery) ToSQL() (string, []interface{}, error) {
	// key := fmt.Sprintf("%s_%s_%s", q.tenantDb.DbName, q.tableName, q.entityType.String())

	if q.Err != nil {
		return "", nil, q.Err
	}
	args := []interface{}{}
	fields := []string{}

	for fieldNAme, arg := range q.fieldMap {
		args = append(args, arg)
		fields = append(fields, utils.ToSnakeCase(fieldNAme))

	}
	autoKeyFieldName := ""
	if q.autoKey != nil {
		autoKeyFieldName = q.autoKey.FieldName
	}
	// sql := ""
	// if v, ok := cacheInsertQuery.Load(key); ok {
	// 	sql = v.(string)
	// } else {

	// }
	sql := q.dialect.BuildSqlInsert(q.tableName, autoKeyFieldName, fields...)
	// cacheInsertQuery.Store(key, sql)
	return sql, args, nil
}
func (q *InsertQuery) ExecWithContext(context context.Context) (interface{}, error) {
	sql, args, err := q.ToSQL()
	if err != nil {
		return nil, err
	}
	// start := time.Now()

	row := q.tenantDb.DB.QueryRowContext(context, sql, args...)
	if row.Err() != nil {
		return nil, row.Err()
	}
	if q.autoKeyType != nil {
		retId := reflect.New(q.autoKeyType)
		err = row.Scan(retId.Interface())
		if err != nil {
			return nil, err
		}
		valueOfqData := reflect.ValueOf(q.data)
		valueOfqData.Elem().FieldByName(q.autoKey.fieldTag.Field.Name).Set(retId.Elem())

	}
	// n := time.Since(start).Microseconds()
	// fmt.Printf("exec time %d us\n", n)

	return q.data, nil

}
func (q *InsertQuery) Exec() (interface{}, error) {
	sql, args, err := q.ToSQL()
	if err != nil {
		return nil, err
	}
	// start := time.Now()
	row := q.tenantDb.DB.QueryRow(sql, args...)
	if row.Err() != nil {
		return nil, row.Err()
	}
	if q.autoKeyType != nil {
		retId := reflect.New(q.autoKeyType)
		err = row.Scan(retId.Interface())
		// n := time.Since(start).Nanoseconds()
		// fmt.Printf("exec sql time :\t\t%d \t\tnns\n", n)
		if err != nil {
			return nil, err
		}
		valueOfqData := reflect.ValueOf(q.data)
		valueOfqData.Elem().FieldByName(q.autoKey.fieldTag.Field.Name).Set(retId.Elem())

	}

	return q.data, nil

}
