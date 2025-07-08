package expr

import (
	"fmt"
	"strings"
	"unvs-orm/internal"

	"github.com/xwb1989/sqlparser"
)

type expressionCompileResult struct {
	Syntax  string
	Context *map[string]string
}

func (c *expression) addTables(tables *[]string, context *map[string]string, tableNames ...string) {
	for _, tableName := range tableNames {
		if _, ok := (*context)[tableName]; !ok { // Đọc từ map thông thường
			(*tables) = append((*tables), tableName)
			(*context)[tableName] = "T" + fmt.Sprintf("%d", len((*context))+1)
			// Ghi vào map thông thường
		}
	}
}
func (e *expression) compile(expr interface{}, tables *[]string, context *map[string]string, isFunctionParamCompiler bool, requireAlias bool) (*expressionCompileResult, error) {

	switch expr := expr.(type) {
	case *sqlparser.SQLVal:
		return e.compileSQLVal(expr)

	case *sqlparser.ColName:
		tableNameFromSyntax := expr.Qualifier.ToViewName().Name.CompliantName()
		tableName := tableNameFromSyntax
		if dbTableName := internal.Utils.GetDbTableName(tableNameFromSyntax); dbTableName != "" {
			tableName = dbTableName
		}
		if tableName == "" {
			return nil, fmt.Errorf("table name not found for %s", expr)
		}
		fieldName := expr.Name.CompliantName()
		metaInfo := internal.Utils.GetMetaInfoByTableName(tableName)

		if metaInfo != nil {
			if _, ok := metaInfo[strings.ToLower(fieldName)]; !ok {
				return nil, fmt.Errorf("field %s not found in table %s", fieldName, tableName)
			} else {
				fieldName = internal.Utils.ToSnakeCase(fieldName)
			}
		}
		if !isFunctionParamCompiler {
			e.addTables(tables, context, tableName)
			tableAlias := (*context)[tableName]
			ret := e.Quote(tableAlias, fieldName) + " AS " + e.Quote(expr.Name.CompliantName())
			return &expressionCompileResult{Syntax: ret, Context: context}, nil
		} else {
			e.addTables(tables, context, tableName)
			tableAlias := (*context)[tableName]
			return &expressionCompileResult{Syntax: e.Quote(tableAlias, fieldName), Context: context}, nil
		}

	case *sqlparser.AliasedExpr:
		ret, err := e.compile(expr.Expr, tables, context, true, requireAlias)
		if err != nil {
			return nil, err
		}
		if isFunctionParamCompiler {
			return ret, nil
		} else {
			if expr.As.CompliantName() != "" {
				synTax := ret.Syntax + " AS " + e.Quote(expr.As.CompliantName())

				return &expressionCompileResult{Syntax: synTax}, nil
			} else {
				return ret, nil
			}
		}
	case *sqlparser.FuncExpr:
		return e.funcExpr(expr, tables, context, requireAlias)
	case sqlparser.SelectExprs:

		retStr := []string{}

		for _, expr := range expr {
			ret, err := e.compile(expr, tables, context, isFunctionParamCompiler, requireAlias)
			if err != nil {
				return nil, err
			}

			retStr = append(retStr, ret.Syntax)
		}
		return &expressionCompileResult{Syntax: strings.Join(retStr, ", ")}, nil

	case *sqlparser.SelectExprs:
		return e.compile(&expr, tables, context, isFunctionParamCompiler, requireAlias)
	case sqlparser.BinaryExpr:
		return e.compileBinaryExpr(&expr, tables, context, false, requireAlias)
	case *sqlparser.BinaryExpr:
		return e.compileBinaryExpr(expr, tables, context, false, requireAlias)
	case *sqlparser.AndExpr:
		return e.AndExpr(expr, tables, context, isFunctionParamCompiler, requireAlias)
	case sqlparser.AndExpr:
		return e.AndExpr(&expr, tables, context, isFunctionParamCompiler, requireAlias)
	case *sqlparser.OrExpr:
		return e.OrExpr(expr, tables, context, isFunctionParamCompiler, requireAlias)
	case sqlparser.OrExpr:
		return e.OrExpr(&expr, tables, context, isFunctionParamCompiler, requireAlias)
	case *sqlparser.ComparisonExpr:
		return e.ComparisonExpr(expr, tables, context, isFunctionParamCompiler, requireAlias)
	case sqlparser.ComparisonExpr:
		return e.ComparisonExpr(&expr, tables, context, isFunctionParamCompiler, requireAlias)
	default:
		return nil, fmt.Errorf("not support %s in file orm/expression.getField.go", expr)
	}

}
