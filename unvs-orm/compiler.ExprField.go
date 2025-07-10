package orm

import (
	EXPR "unvs-orm/expr"
)

func (c *CompilerUtils) resolveExprField(tables *[]string, context *map[string]string, f *exprField, extractAlias, applyContext bool) (*resolverResult, error) {
	if bf, ok := f.UnderField.(*exprField); ok {
		selectCmp := EXPR.NewExpressionCompiler(c.dialect.driverName())
		r, err := selectCmp.Compile("", tables, context, bf.Stmt, false, true) //<-- when do compiling use alias from context
		if err != nil {
			return nil, err
		}
		return &resolverResult{
			Syntax: r.Syntax,
			Args:   r.Args,
		}, nil
	} else {
		return c.Resolve(tables, context, f.UnderField, extractAlias, applyContext)
	}
}
