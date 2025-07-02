package unvsef

import (
	"reflect"
)

type Query struct {
	table    interface{}
	tenantDb *TenantDb
}

func From(table interface{}) *Query {
	return &Query{
		table: table,
	}
}

func LeftJoin(table interface{}) interface{} {
	return &joinExpr{
		JoinType: "LEFT JOIN",
		Table:    table,
	}
}

func RightJoin(table interface{}) interface{} {
	return &joinExpr{
		JoinType: "RIGHT JOIN",
		Table:    table,
	}
}

func InnerJoin(joinExpr *BinaryField) interface{} {
	panic("not implemented")
}

func CrossJoin(table interface{}) interface{} {
	return &joinExpr{
		JoinType: "CROSS JOIN",
		Table:    table,
	}
}

type joinExpr struct {
	JoinType string
	Table    interface{}
	On       interface{} // can be nil for CROSS JOIN or delayed assignment
}

type SelectQuery struct {
	owner    *Query
	table    interface{}
	Fields   []interface{}
	WhereEx  interface{}
	Group    []interface{}
	HavingEx interface{}
	Order    []interface{}
}

type SelectQueryWithOrder struct {
	*SelectQuery
	LimitEx  *int
	OffsetEx *int
}

func (q *SelectQueryWithOrder) String() string {
	sql, _ := q.ToSQL(q.owner.tenantDb.Dialect)

	return sql
}
func (q *Query) convertToAlias(field interface{}) interface{} {
	if _, ok := field.(*AliasField); !ok {
		return field
	}
	panic("not implemented")

}
func (q *Query) Select(fields ...interface{}) *SelectQuery {
	selectFields := make([]interface{}, 0, len(fields))
	for _, f := range fields {
		// selectFields = append(selectFields, q.convertToAlias(f))
		if _, ok := f.(*AliasField); !ok {
			valType := reflect.ValueOf(f)
			if valType.Kind() == reflect.Ptr {
				valType = valType.Elem()
			}
			aliasField := valType.FieldByName("FieldName")

			if aliasField.IsValid() {
				fieldName := aliasField.String()
				selectFields = append(selectFields, &AliasField{f, fieldName})
			}

		} else {
			selectFields = append(selectFields, f.(*AliasField))
		}
	}
	// for _, f := range fields {
	// 	dbField := reflect.ValueOf(f).FieldByName("DbField")
	// 	if dbField.Kind() == reflect.Ptr {
	// 		dbField = dbField.Elem()
	// 	}
	// 	if dbField.IsValid() {
	// 		aliasField := dbField.FieldByName("Alias")
	// 		fieldName := dbField.FieldByName("FieldName")
	// 		strFieldName := fieldName.String()
	// 		if aliasField.IsValid() {
	// 			aliasField.SetString(strFieldName)
	// 		}

	// 	}
	// }
	return &SelectQuery{
		owner:  q,
		table:  q.table,
		Fields: selectFields,
	}
}

func (q *SelectQuery) Where(where interface{}) *SelectQuery {
	q.WhereEx = where
	return q
}

func (q *SelectQuery) GroupBy(fields ...interface{}) *SelectQuery {
	q.Group = fields
	return q
}

func (q *SelectQuery) Having(where interface{}) *SelectQuery {
	q.HavingEx = where
	return q
}

func (q *SelectQuery) OrderBy(fields ...interface{}) *SelectQueryWithOrder {
	q.Order = fields
	return &SelectQueryWithOrder{
		SelectQuery: q,
	}
}

func (q *SelectQueryWithOrder) Limit(limit int) *SelectQueryWithOrder {
	q.LimitEx = &limit
	return q
}

func (q *SelectQueryWithOrder) Offset(offset int) *SelectQueryWithOrder {
	q.OffsetEx = &offset
	return q
}

func (q *SelectQuery) ToSQL(d Dialect) (string, []interface{}) {
	return (&SelectQueryWithOrder{SelectQuery: q}).ToSQL(d)
}
func (q *SelectQuery) execToByType(typ reflect.Type) (interface{}, error) {
	sqlStmt, args := q.ToSQL(q.owner.tenantDb.Dialect)
	rows, err := q.owner.tenantDb.DB.Query(sqlStmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return utils.fetchAllRows(rows, typ)
}
func (q *SelectQuery) ExecTo(entity interface{}) (interface{}, error) {
	typ := reflect.TypeOf(entity)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() == reflect.Slice {
		typ = typ.Elem()
	}
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return q.execToByType(typ)

}

func (q *SelectQueryWithOrder) ToSQL(d Dialect) (string, []interface{}) {
	sql := "SELECT"
	args := []interface{}{}

	// Fields
	if len(q.Fields) == 0 {
		sql += " *"
	} else {
		fields := make([]string, len(q.Fields))
		for i, f := range q.Fields {

			expr, a := compiler.Compile(f, d)
			fields[i] = expr
			args = append(args, a...)
		}
		sql += " " + utils.Join(fields, ", ")
	}

	// FROM
	switch jt := q.table.(type) {
	case *joinExpr:
		tableSQL, tableArgs := compiler.Compile(jt.Table, d)
		sql += " FROM " + tableSQL + " " + jt.JoinType
		args = append(args, tableArgs...)
		if jt.JoinType != "CROSS JOIN" && jt.On != nil {
			onSQL, onArgs := compiler.Compile(jt.On, d)
			sql += " ON " + onSQL
			args = append(args, onArgs...)
		}
	default:
		tableExpr, tableArgs := compiler.Compile(q.table, d)
		sql += " FROM " + tableExpr
		args = append(args, tableArgs...)
	}

	// WHERE
	if q.WhereEx != nil {
		sql += " WHERE "
		expr, a := compiler.Compile(q.WhereEx, d)
		sql += expr
		args = append(args, a...)
	}

	// GROUP BY
	if len(q.Group) > 0 {
		groups := make([]string, len(q.Group))
		for i, g := range q.Group {
			expr, a := compiler.Compile(g, d)
			groups[i] = expr
			args = append(args, a...)
		}
		sql += " GROUP BY " + utils.Join(groups, ", ")
	}

	// HAVING
	if q.HavingEx != nil {
		sql += " HAVING "
		expr, a := compiler.Compile(q.HavingEx, d)
		sql += expr
		args = append(args, a...)
	}

	// ORDER BY
	if len(q.Order) > 0 {
		orders := make([]string, len(q.Order))
		for i, o := range q.Order {
			expr, a := compiler.Compile(o, d)
			orders[i] = expr
			args = append(args, a...)
		}
		sql += " ORDER BY " + utils.Join(orders, ", ")
	}

	// LIMIT & OFFSET
	limOff := d.MakeLimitOffset(q.LimitEx, q.OffsetEx)
	if limOff != "" {
		sql += " " + limOff
		// if q.LimitEx != nil {
		// 	args = append(args, *q.LimitEx)
		// }
		// if q.OffsetEx != nil {
		// 	args = append(args, *q.OffsetEx)
		// }
	}

	return utils.replacePlaceHolder(d.GetParamPlaceholder(), sql), args
}
func (q *SelectQueryWithOrder) execToByType(typ reflect.Type) (interface{}, error) {
	sqlStmt, args := q.ToSQL(q.owner.tenantDb.Dialect)
	rows, err := q.owner.tenantDb.DB.Query(sqlStmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return utils.fetchAllRows(rows, typ)

}
