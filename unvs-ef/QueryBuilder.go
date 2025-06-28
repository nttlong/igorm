package unvsef

type Query struct {
	table interface{}
}

func From(table interface{}) *Query {
	return &Query{
		table: table,
	}
}

func LeftJoin(table interface{}) interface{} {
	return &JoinExpr{
		JoinType: "LEFT JOIN",
		Table:    table,
	}
}

func RightJoin(table interface{}) interface{} {
	return &JoinExpr{
		JoinType: "RIGHT JOIN",
		Table:    table,
	}
}

func InnerJoin(table interface{}) interface{} {
	return &JoinExpr{
		JoinType: "INNER JOIN",
		Table:    table,
	}
}

func CrossJoin(table interface{}) interface{} {
	return &JoinExpr{
		JoinType: "CROSS JOIN",
		Table:    table,
	}
}

type JoinExpr struct {
	JoinType string
	Table    interface{}
	On       interface{} // can be nil for CROSS JOIN or delayed assignment
}

type SelectQuery struct {
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

func (q *Query) Select(fields ...interface{}) *SelectQuery {
	return &SelectQuery{
		table:  q.table,
		Fields: fields,
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
	case *JoinExpr:
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
		if q.LimitEx != nil {
			args = append(args, *q.LimitEx)
		}
		if q.OffsetEx != nil {
			args = append(args, *q.OffsetEx)
		}
	}

	return sql, args
}
