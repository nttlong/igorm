package dbv

import (
	// EXPR "dbv/expr"
	"dbv/tenantDB"
	"strconv"
	"strings"
	"sync"
)

type QueryParts struct {
	selectFields []string
	argsSelect   []interface{}
	fromExpr     string
	whereExprs   string
	whereArgs    []interface{}
	orderByExprs []string
	orderByArgs  []interface{}
	groupByExprs []string
	groupByArgs  []interface{}
	Limit        *int
	Offset       *int
	Err          error
	having       string
	havingArgs   []interface{}
}

func NewQueryParts() *QueryParts {
	return &QueryParts{}
}

type selectExpr struct {
	Expr string
	Args []interface{}
}
type litOfString struct {
	val string
}

func Lit(str string) litOfString {
	return litOfString{val: str}
}

/*
Example:

	Select("id", "name", "age") <-- no args

	Select("id", "name", "age", 123) <-- with args

	In case arg is string select can not recognize it as parameter and will be treated as string literal
	Select("concat(firstName, ?,lastName)",Lit(" ")) <-- with args
*/
func (q *QueryParts) Select(exprsAndArgs ...interface{}) *QueryParts {
	// if q.selectFields == nil {
	// 	q.selectFields = []selectExpr{}
	// }

	// if len(exprsAndArgs) == 0 {
	// 	return q
	// }

	// Tách phần expression và phần args
	expressions := []string{}
	args := []interface{}{}

	for _, item := range exprsAndArgs {
		switch v := item.(type) {
		case string:
			expressions = append(expressions, v)
		case litOfString:
			// LitOfString là literal, sẽ dùng làm arg
			args = append(args, v)
		default:
			// Các kiểu khác cũng là arg
			args = append(args, v)
		}
	}
	q.selectFields = expressions
	q.argsSelect = args

	return q
}
func (q *QueryParts) InnerJoin(table string, onExpr string, args ...interface{}) *QueryParts {
	q.fromExpr = q.fromExpr + " INNER JOIN " + table + " ON " + onExpr
	q.argsSelect = append(q.argsSelect, args...)
	return q
}
func (q *QueryParts) LeftJoin(table string, onExpr string, args ...interface{}) *QueryParts {
	q.fromExpr = q.fromExpr + " LEFT JOIN " + table + " ON " + onExpr
	q.argsSelect = append(q.argsSelect, args...)
	return q
}
func (q *QueryParts) From(table string) *QueryParts {
	q.fromExpr = table
	return q
}

func (q *QueryParts) Where(expr string, args ...interface{}) *QueryParts {
	q.whereExprs = expr
	q.whereArgs = args
	return q
}

func (q *QueryParts) OrderBy(args ...interface{}) *QueryParts {
	exprs := []string{}
	for _, item := range args {
		switch v := item.(type) {
		case string:
			exprs = append(exprs, v)
		case litOfString:
			// LitOfString là literal, sẽ dùng làm arg
			q.orderByArgs = append(q.orderByArgs, v)
		default:
			// Các kiểu khác cũng là arg
			q.orderByArgs = append(q.orderByArgs, v)
		}
	}
	q.orderByExprs = exprs
	return q
}

func (q *QueryParts) GroupBy(args ...interface{}) *QueryParts {
	exprs := []string{}
	for _, item := range args {
		switch v := item.(type) {
		case string:
			exprs = append(exprs, v)
		case litOfString:
			// LitOfString là literal, sẽ dùng làm arg
			q.groupByArgs = append(q.groupByArgs, v)
		default:
			// Các kiểu khác cũng là arg
			q.groupByArgs = append(q.groupByArgs, v)
		}
	}
	q.groupByExprs = exprs
	return q
}

func (q *QueryParts) OffsetLimit(offset, limit int) *QueryParts {
	q.Limit = &limit
	q.Offset = &offset
	return q
}

type initBuildSQL struct {
	once sync.Once
	val  string
	err  error
}

var cacheBuildSQL sync.Map

func (q *QueryParts) buildSQL(db *tenantDB.TenantDB) (string, error) {
	var sb strings.Builder

	compiler, err := CompileJoin(q.fromExpr, db)
	// if !strings.Contains(q.fromExpr, " JOIN ") {
	// 	q.fromExpr = q.fromExpr + " AS t"

	// }
	if err != nil {

		return "", err
	} else {
		q.fromExpr = compiler.content
	}
	// SELECT
	sb.WriteString("SELECT ")
	if len(q.selectFields) == 0 {
		sb.WriteString("*")
	} else {

		selectSyntax := strings.Join(q.selectFields, ", ")
		err = compiler.buildSelectField(selectSyntax)
		if err != nil {
			q.Err = err
			return "", nil
		}
		sb.WriteString(compiler.content)
	}

	// FROM
	if q.fromExpr != "" {
		sb.WriteString(" FROM ")
		sb.WriteString(q.fromExpr)
	}

	// WHERE
	if len(q.whereExprs) > 0 {
		sb.WriteString(" WHERE ")
		compiler.context.purpose = build_purpose_where
		err = compiler.buildWhere(q.whereExprs)
		if err != nil {
			return "", nil
		}

		sb.WriteString(compiler.content)
	}

	// GROUP BY
	if len(q.groupByExprs) > 0 {
		sb.WriteString(" GROUP BY ")

		groupSyntax := strings.Join(q.groupByExprs, ", ")
		err = compiler.buildSelectField(groupSyntax)
		if err != nil {
			return "", nil
		}
		sb.WriteString(groupSyntax)

	}

	// ORDER BY
	if len(q.orderByExprs) > 0 {
		sb.WriteString(" ORDER BY ")
		sortExpr := strings.Join(q.orderByExprs, ", ")
		err = compiler.buildSortField(sortExpr)
		if err != nil {
			return "", nil
		}
		sb.WriteString(compiler.content)
	}
	return sb.String(), nil

}
func (q *QueryParts) BuildSQL(db *tenantDB.TenantDB) (string, []interface{}) {
	key := db.GetDriverName() + "://" + db.GetDBName() + "/" + q.fromExpr + "/" + q.whereExprs
	key += "/" + strings.Join(q.orderByExprs, ",")
	key += "/" + strings.Join(q.groupByExprs, ",")
	key += "/" + strings.Join(q.selectFields, ",")
	key += "/" + q.having
	actual, _ := cacheBuildSQL.LoadOrStore(key, &initBuildSQL{})
	initBuild := actual.(*initBuildSQL)
	initBuild.once.Do(func() {
		sql, err := q.buildSQL(db)
		initBuild.val = sql
		initBuild.err = err
	})
	if initBuild.err != nil {
		q.Err = initBuild.err
		return "", nil
	}
	reqSql := initBuild.val
	if db.GetDriverName() == "sqlserver" {
		//OFFSET 100 ROWS FETCH NEXT 0 ROWS ONLY;
		if q.Limit == nil {
			*q.Limit = 0
		}
		if q.Offset == nil {
			*q.Offset = 0
		}
		reqSql += " OFFSET " + strconv.Itoa(*q.Offset) + " ROWS FETCH NEXT " + strconv.Itoa(*q.Limit) + " ROWS ONLY"

	} else {
		// LIMIT OFFSET
		if q.Limit != nil {
			reqSql += " LIMIT " + strconv.Itoa(*q.Limit)
			// sb.WriteString(" LIMIT ?")

		}
		if q.Offset != nil {
			//sb.WriteString(" OFFSET ?")
			reqSql += " OFFSET " + strconv.Itoa(*q.Offset)

		}
	}
	args := []interface{}{}

	args = append(args, q.argsSelect...)
	args = append(args, q.whereArgs...)
	if q.Limit != nil {
		args = append(args, *q.Limit)
	}
	if q.Offset != nil {

		args = append(args, *q.Offset)
	}
	for i, a := range args {
		if litOfString, isLit := a.(litOfString); isLit {
			args[i] = litOfString.val
		}
	}
	return reqSql, args
}
func Qr() *QueryParts {
	return NewQueryParts()
}
func init() {
	tenantDB.OnNewQrFn = func() interface{} {
		return NewQueryParts()
	}
	tenantDB.OnQrSelect = func(q interface{}, exprsAndArgs ...interface{}) interface{} {
		qr := q.(*QueryParts)
		return qr.Select(exprsAndArgs...)
	}
	tenantDB.OnInnerJoin = func(q interface{}, table string, onExpr string, args ...interface{}) interface{} {
		qr := q.(*QueryParts)
		return qr.InnerJoin(table, onExpr, args...)
	}
	tenantDB.OnLeftJoin = func(q interface{}, table string, onExpr string, args ...interface{}) interface{} {
		qr := q.(*QueryParts)
		return qr.LeftJoin(table, onExpr, args...)
	}
	tenantDB.OnWhere = func(q interface{}, expr string, args ...interface{}) interface{} {
		qr := q.(*QueryParts)
		return qr.Where(expr, args...)
	}
	tenantDB.OnOrderBy = func(q interface{}, args ...interface{}) interface{} {
		qr := q.(*QueryParts)
		return qr.OrderBy(args...)
	}
	tenantDB.OnGroupBy = func(q interface{}, args ...interface{}) interface{} {
		qr := q.(*QueryParts)
		return qr.GroupBy(args...)
	}
	tenantDB.OnOffsetLimit = func(q interface{}, offset, limit int) interface{} {
		qr := q.(*QueryParts)
		return qr.OffsetLimit(offset, limit)
	}
	tenantDB.OnBuildSql = func(q interface{}, db *tenantDB.TenantDB) (string, []interface{}) {
		qr := q.(*QueryParts)
		return qr.BuildSQL(db)
	}
	tenantDB.OnFrom = func(q interface{}, table string, args ...interface{}) interface{} {
		qr := q.(*QueryParts)
		return qr.From(table)
	}
}
