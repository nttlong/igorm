package vdb

import (
	"strconv"
	"vdb/sqlparser"
)

func (compiler *exprReceiver) SQLVal(context *exprCompileContext, expr *sqlparser.SQLVal) (string, error) {
	switch expr.Type {
	case sqlparser.StrVal:
		return context.dialect.ToText(string(expr.Val)), nil
	case sqlparser.IntVal:
		return string(expr.Val), nil
	case sqlparser.FloatVal:
		return string(expr.Val), nil
	case sqlparser.ValArg:
		if context.paramIndex == 0 {
			context.paramIndex = 1
		}

		strIndex := string(expr.Val[2:len(expr.Val)])
		if _, err := strconv.Atoi(strIndex); err == nil {
			defer func() {
				context.paramIndex++
			}()
			return context.dialect.ToParam(context.paramIndex), nil
		} else {
			return string(expr.Val), nil
		}

	}

	return string(expr.Val), nil

}
