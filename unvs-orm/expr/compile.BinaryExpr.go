package expr

import (
	"fmt"

	"github.com/xwb1989/sqlparser"
)

func (e *expression) compileBinaryExpr(expr *sqlparser.BinaryExpr, tables *[]string, context *map[string]string, isFunctionParamCompiler, extractAlias, applyContext bool) (*expressionCompileResult, error) {
	left, err := e.compile(expr.Left, tables, context, true, extractAlias, applyContext)

	if err != nil {
		return nil, err
	}

	right, err := e.compile(expr.Right, tables, context, true, extractAlias, applyContext)
	if err != nil {
		return nil, err
	}

	ret := fmt.Sprintf("%s %s %s", left.Syntax, expr.Operator, right.Syntax)
	return &expressionCompileResult{Syntax: ret}, nil
}
