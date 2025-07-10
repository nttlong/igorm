package expr

import (
	"strings"

	"github.com/xwb1989/sqlparser"
)

func (e *expression) AndExpr(expr *sqlparser.AndExpr, tables *[]string, context *map[string]string, isFunctionParamCompiler, extractAlias, applyContext bool) (*expressionCompileResult, error) {
	left, err := e.compile(expr.Left, tables, context, isFunctionParamCompiler, extractAlias, applyContext)
	if err != nil {
		return nil, err
	}

	right, err := e.compile(expr.Right, tables, context, isFunctionParamCompiler, extractAlias, applyContext)
	if err != nil {
		return nil, err
	}

	ret := []string{left.Syntax + " AND " + right.Syntax}
	return &expressionCompileResult{Syntax: strings.Join(ret, " AND ")}, nil
}
func (e *expression) OrExpr(expr *sqlparser.OrExpr, tables *[]string, context *map[string]string, isFunctionParamCompiler, extractAlias, applyContext bool) (*expressionCompileResult, error) {
	left, err := e.compile(expr.Left, tables, context, isFunctionParamCompiler, extractAlias, applyContext)
	if err != nil {
		return nil, err
	}

	right, err := e.compile(expr.Right, tables, context, isFunctionParamCompiler, extractAlias, applyContext)
	if err != nil {
		return nil, err
	}

	ret := []string{left.Syntax + " OR " + right.Syntax}
	return &expressionCompileResult{Syntax: strings.Join(ret, " OR ")}, nil
}
