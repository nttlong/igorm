package eorm

import "github.com/xwb1989/sqlparser"

func (compiler *exprReceiver) AliasedTableExpr(context *exprCompileContext, expr *sqlparser.AliasedTableExpr) (string, error) {

	tableName := expr.As.CompliantName()
	if tableName == "$$$$$$$$$$$$$$" {
		return "", nil
	}
	if tableName == "" {
		return compiler.compile(context, expr.Expr)
	}
	return tableName, nil

}
