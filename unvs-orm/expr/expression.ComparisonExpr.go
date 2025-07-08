package expr

import "github.com/xwb1989/sqlparser"

func (e *expression) ComparisonExpr(expr *sqlparser.ComparisonExpr, tables *[]string, context *map[string]string, isFunctionParamCompiler bool, requireAlias bool) (*expressionCompileResult, error) {
	// TODO: implement this function
	left, err := e.compile(expr.Left, tables, context, isFunctionParamCompiler, requireAlias)
	if err != nil {
		return nil, err
	}

	right, err := e.compile(expr.Right, tables, context, isFunctionParamCompiler, requireAlias)
	if err != nil {
		return nil, err
	}

	syntax := left.Syntax + " " + string(expr.Operator) + " " + right.Syntax
	return &expressionCompileResult{Syntax: syntax}, nil

}
