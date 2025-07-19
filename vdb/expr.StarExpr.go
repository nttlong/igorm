package vdb

import "vdb/sqlparser"

func (compiler *exprReceiver) StarExpr(context *exprCompileContext, expr *sqlparser.StarExpr) (string, error) {
	if expr.TableName.IsEmpty() {
		return "*", nil
	}
	panic("not implemented")

}
