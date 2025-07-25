package vdb

import (
	"fmt"

	"vdb/sqlparser"
)

func (compiler *exprReceiver) JoinCondition(context *exprCompileContext, expr *sqlparser.JoinCondition) (string, error) {
	if sqlparserComparisonExpr, ok := expr.On.(*sqlparser.ComparisonExpr); ok {
		return compiler.ComparisonExpr(context, *sqlparserComparisonExpr)
	}
	panic(fmt.Errorf("unsupported expression type %T in file eorm/expr.go, line 17", expr))

}
