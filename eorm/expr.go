package eorm

import (
	"fmt"

	"github.com/xwb1989/sqlparser"
)

type exprReceiver struct {
}

func (e *exprReceiver) compile(context *exprCompileContext, expr interface{}) (string, error) {
	switch expr := expr.(type) {
	case *sqlparser.ComparisonExpr:
		return e.ComparisonExpr(context, *expr)
	case *sqlparser.ColName:
		return e.ColName(context, *expr)
	case *sqlparser.AndExpr:
		return e.AndExpr(context, expr)
	case *sqlparser.SQLVal:
		return e.SQLVal(context, expr)
	case *sqlparser.StarExpr:
		return e.StarExpr(context, expr)
	case *sqlparser.JoinTableExpr:
		return e.JoinTableExpr(context, expr)
	case *sqlparser.AliasedTableExpr:
		return e.AliasedTableExpr(context, expr)
	case sqlparser.JoinCondition:
		return e.JoinCondition(context, &expr)
	case sqlparser.TableName:
		return e.TableName(context, &expr)

	default:
		panic(fmt.Errorf("unsupported expression type %T in file eorm/expr.go, line 17", expr))
	}
	return "", nil
}

var exprs = &exprReceiver{}
