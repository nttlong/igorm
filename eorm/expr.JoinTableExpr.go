package eorm

import (
	"fmt"

	"eorm/sqlparser"
)

func (compiler *exprReceiver) JoinTableExpr(context *exprCompileContext, expr *sqlparser.JoinTableExpr) (string, error) {
	left, err := compiler.compile(context, expr.LeftExpr)
	if err != nil {
		return "", err
	}
	right, err := compiler.compile(context, expr.RightExpr)
	if err != nil {
		return "", err
	}
	on, err := compiler.compile(context, expr.Condition)
	if err != nil {
		return "", err
	}
	if expr.Join == "join" {
		if left == "" {
			return fmt.Sprintf("INNER JOIN %s ON %s", right, on), nil
		}
		return fmt.Sprintf("%s INNER JOIN %s ON %s", left, right, on), nil
	}

	panic("not implemented")

}
