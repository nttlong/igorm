package dbv

import (
	"dbv/sqlparser"
	"strconv"
)

func (compiler *exprReceiver) SQLVal(context *exprCompileContext, expr *sqlparser.SQLVal) (string, error) {
	if expr.Type == sqlparser.StrVal {

		return context.dialect.ToText(string(expr.Val)), nil
	} else {
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

}
