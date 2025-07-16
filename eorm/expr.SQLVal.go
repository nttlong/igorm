package dbv

import "dbv/sqlparser"

func (compiler *exprReceiver) SQLVal(context *exprCompileContext, expr *sqlparser.SQLVal) (string, error) {
	if expr.Type == sqlparser.StrVal {

		return context.dialect.ToText(string(expr.Val)), nil
	} else {
		return string(expr.Val), nil
	}

}
