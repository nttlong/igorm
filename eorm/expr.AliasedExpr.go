package eorm

import (
	"github.com/xwb1989/sqlparser"
)

func (compiler *exprReceiver) AliasedExpr(context *exprCompileContext, expr *sqlparser.AliasedExpr) (string, error) {
	return compiler.compile(context, expr.Expr)

}
