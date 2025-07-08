package orm

import (
	EXPR "unvs-orm/expr"
)

func (c *CompilerUtils) resolveExprField(tables *[]string, context *map[string]string, f *exprField, requireAlias bool) (*resolverResult, error) {
	if bf, ok := f.UnderField.(*ExprBase); ok {
		selectCmp := EXPR.NewExpressionCompiler(c.dialect.driverName())
		r, err := selectCmp.Compile(tables, context, bf.Stmt, true)
		if err != nil {
			return nil, err
		}
		return &resolverResult{
			Syntax:       r.Syntax,
			Args:         r.Args,
			buildContext: context,
			Tables:       tables,
		}, nil
	} else {
		return c.Resolve(tables, context, f.UnderField, requireAlias)
	}
}
