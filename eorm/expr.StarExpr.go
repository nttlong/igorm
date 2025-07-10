package eorm

import "github.com/xwb1989/sqlparser"

func (compiler *exprReceiver) StarExpr(context *exprCompileContext, expr *sqlparser.StarExpr) (string, error) {
	if expr.TableName.IsEmpty() {
		return "*", nil
	}
	panic("not implemented")

}
