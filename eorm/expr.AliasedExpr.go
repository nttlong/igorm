package eorm

import (
	"eorm/sqlparser"
)

func (compiler *exprReceiver) AliasedExpr(context *exprCompileContext, expr *sqlparser.AliasedExpr) (string, error) {
	return compiler.compile(context, expr.Expr)

}
