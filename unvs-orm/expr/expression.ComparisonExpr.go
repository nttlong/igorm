package expr

import "github.com/xwb1989/sqlparser"

func (e *expression) ComparisonExpr(expr *sqlparser.ComparisonExpr, tables *[]string, context *map[string]string, isFunctionParamCompiler, extractAlias, applyContext bool) (*expressionCompileResult, error) {
	// TODO: implement this function
	left, err := e.compile(expr.Left, tables, context, isFunctionParamCompiler, extractAlias, applyContext)
	if err != nil {
		return nil, err
	}

	right, err := e.compile(expr.Right, tables, context, isFunctionParamCompiler, extractAlias, applyContext)
	if err != nil {
		return nil, err
	}

	syntax := left.Syntax + " " + string(expr.Operator) + " " + right.Syntax
	return &expressionCompileResult{Syntax: syntax}, nil

}
