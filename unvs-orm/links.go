package orm

import (
	"fmt"
	expression "unvs-orm/expr"
)

var OnGetQuoteFunc = func(dbDriver expression.DB_TYPE, str ...string) string {

	if dbDriver == expression.DB_TYPE_MYSQL {
		panic(fmt.Errorf("not implemented mysql dialect in file orm/links.go, line %d", 12))
	}
	if dbDriver == expression.DB_TYPE_POSTGRES {
		panic(fmt.Errorf("not implemented mysql dialect in file orm/links.go, line %d", 17))
	}
	if dbDriver == expression.DB_TYPE_MSSQL {

		return MssqlCompiler.Quote(str...)
	}

	panic(fmt.Sprintf("not support dialect for %s, file orm/links.go, line %d", dbDriver, 21))

}

func init() {
	expression.OnGetQuoteFunc = OnGetQuoteFunc
	expression.OnCompileFunc = func(dbDriver expression.DB_TYPE, tables *[]string, context *map[string]string, caller interface{}, extractAlias, applyContext bool) (*expression.ResolverResult, error) {

		if dbDriver == expression.DB_TYPE_MSSQL {
			//if exprCall,ok:=caller.(*expression.ExpressionTest);ok{

			result, err := MssqlCompiler.Resolve(tables, context, caller, extractAlias, applyContext)
			if err != nil {
				return nil, err
			}
			return &expression.ResolverResult{
				Syntax: result.Syntax,
				Args:   result.Args,
			}, nil

		}

		panic(fmt.Sprintf("not support dialect for %s, file orm/links.go, line %d", dbDriver, 45))
	}
}
