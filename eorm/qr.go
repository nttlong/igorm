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
	groupByExprs []string
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
	// argIdx := 0
	// for _, exprStr := range expressions {
	// 	numArgs := strings.Count(exprStr, "?")
	// 	if argIdx+numArgs > len(args) {
	// 		panic(fmt.Sprintf("Select: not enough arguments for expression '%s' — expected %d but have %d left", exprStr, numArgs, len(args)-argIdx))
	// 	}

	// 	exprArgs := []interface{}{}
	// 	for _, a := range args[argIdx : argIdx+numArgs] {
	// 		if lit, isLit := a.(litOfString); isLit {
	// 			exprArgs = append(exprArgs, lit.val)
	// 		} else {
	// 			exprArgs = append(exprArgs, a)
	// 		}
	// 	}

	// 	argIdx += numArgs
	// }

	// if argIdx != len(args) {
	// 	panic(fmt.Sprintf("Select: too many arguments — used %d but received %d", argIdx, len(args)))
	// }

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

func (q *QueryParts) OrderBy(exprs ...string) *QueryParts {
	q.orderByExprs = append(q.orderByExprs, exprs...)
	return q
}

func (q *QueryParts) GroupBy(exprs ...string) *QueryParts {
	q.groupByExprs = append(q.groupByExprs, exprs...)
	return q
}

func (q *QueryParts) LimitOffset(limit, offset int) *QueryParts {
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
		sb.WriteString(strings.Join(q.orderByExprs, ", "))
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

	// LIMIT OFFSET
	if q.Limit != nil {
		reqSql += " LIMIT " + strconv.Itoa(*q.Limit)
		// sb.WriteString(" LIMIT ?")

	}
	if q.Offset != nil {
		//sb.WriteString(" OFFSET ?")
		reqSql += " OFFSET " + strconv.Itoa(*q.Limit)

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
