package orm

import (
	"fmt"
	"strings"

	"github.com/xwb1989/sqlparser"
)

func (e *expression) compile(expr interface{}, isFunctionParamCompiler bool) ([]string, error) {

	switch expr := expr.(type) {
	case *sqlparser.SQLVal:
		return e.compileSQLVal(expr)

	case *sqlparser.ColName:
		tableNameFromSyntax := expr.Qualifier.ToViewName().Name.CompliantName()
		tableName := tableNameFromSyntax
		if dbTableName := Utils.GetDbTableName(tableNameFromSyntax); dbTableName != "" {
			tableName = dbTableName
		}
		if tableName == "" {
			return nil, fmt.Errorf("table name not found for %s", expr)
		}
		fieldName := expr.Name.CompliantName()
		metaInfo := Utils.GetMetaInfoByTableName(tableName)
		if metaInfo != nil {
			if _, ok := metaInfo[strings.ToLower(fieldName)]; !ok {
				return nil, fmt.Errorf("field %s not found in table %s", fieldName, tableName)
			} else {
				fieldName = Utils.ToSnakeCase(fieldName)
			}
		}
		if !isFunctionParamCompiler {
			ret := e.cmp.Quote(tableName, fieldName) + " AS " + e.cmp.Quote(expr.Name.CompliantName())
			return []string{ret}, nil
		} else {
			return []string{e.cmp.Quote(tableName, fieldName)}, nil
		}

	case *sqlparser.AliasedExpr:
		return e.compile(expr.Expr, isFunctionParamCompiler)
	case *sqlparser.FuncExpr:
		return e.funcExpr(expr)
	case sqlparser.SelectExprs:

		retStr := []string{}

		for _, expr := range expr {
			ret, err := e.compile(expr, isFunctionParamCompiler)
			if err != nil {
				return nil, err
			}
			retStr = append(retStr, ret...)
		}
		return retStr, nil

	case *sqlparser.SelectExprs:
		return e.compile(expr, isFunctionParamCompiler)
	default:
		return nil, fmt.Errorf("not support %s", expr)
	}

}
