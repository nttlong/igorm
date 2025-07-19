package vdb

import (
	"vdb/sqlparser"
)

func (compiler *exprReceiver) AliasedExpr(context *exprCompileContext, expr *sqlparser.AliasedExpr) (string, error) {
	return compiler.compile(context, expr.Expr)

}
