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

/*
@extractAlias: true: neu tra cuu trong context ma kg co tim thay bang tao alias moi

	dong thoi insert vao context va tables bieu thuc sinh ra pahi co table alias thay vi ten bang goc,
	false: khong tac dong vao context va tables

	@applyContext: true: áp dụng context cho tên bảng, false: không áp dụng context cho tên bảng
*/
func (e *expression) compile(expr interface{}, tables *[]string, context *map[string]string, isFunctionParamCompiler bool, extractAlias bool, applyContext bool) (*expressionCompileResult, error) {

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
		if !isFunctionParamCompiler { //<-- if not in function

			if extractAlias { //<-- if extract alias from syntax
				e.addTables(tables, context, tableName)
				tableAlias := (*context)[tableName]
				ret := e.Quote(tableAlias, fieldName) + " AS " + e.Quote(expr.Name.CompliantName())
				return &expressionCompileResult{Syntax: ret, Context: context}, nil
			} else if applyContext { //<-- if apply context to syntax use table alias in context
				tableAlias := (*context)[tableName]
				ret := e.Quote(tableAlias, fieldName) + " AS " + e.Quote(expr.Name.CompliantName())
				return &expressionCompileResult{Syntax: ret, Context: context}, nil
			} else { //<-- if not apply context to syntax
				ret := e.Quote(tableName, fieldName) + " AS " + e.Quote(expr.Name.CompliantName())
				return &expressionCompileResult{Syntax: ret, Context: context}, nil
			}
		} else { // in function skip field alias
			if extractAlias { //<-- if extract alias from syntax
				e.addTables(tables, context, tableName)
				tableAlias := (*context)[tableName]
				ret := e.Quote(tableAlias, fieldName)
				return &expressionCompileResult{Syntax: ret, Context: context}, nil
			} else if applyContext { //<-- if apply context to syntax use table alias in context
				tableAlias := (*context)[tableName]
				ret := e.Quote(tableAlias, fieldName)
				return &expressionCompileResult{Syntax: ret, Context: context}, nil
			} else { //<-- if not apply context to syntax
				ret := e.Quote(tableName, fieldName)
				return &expressionCompileResult{Syntax: ret, Context: context}, nil
			}
		}

	case *sqlparser.AliasedExpr:
		ret, err := e.compile(expr.Expr, tables, context, true, extractAlias, applyContext)
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
		return e.funcExpr(expr, tables, context, extractAlias, applyContext)
	case sqlparser.SelectExprs:

		retStr := []string{}

		for _, expr := range expr {
			ret, err := e.compile(expr, tables, context, isFunctionParamCompiler, extractAlias, applyContext)
			if err != nil {
				return nil, err
			}

			retStr = append(retStr, ret.Syntax)
		}
		return &expressionCompileResult{Syntax: strings.Join(retStr, ", ")}, nil

	case *sqlparser.SelectExprs:
		return e.compile(&expr, tables, context, isFunctionParamCompiler, extractAlias, applyContext)
	case sqlparser.BinaryExpr:
		return e.compileBinaryExpr(&expr, tables, context, false, extractAlias, applyContext)
	case *sqlparser.BinaryExpr:
		return e.compileBinaryExpr(expr, tables, context, false, extractAlias, applyContext)
	case *sqlparser.AndExpr:
		return e.AndExpr(expr, tables, context, isFunctionParamCompiler, extractAlias, applyContext)
	case sqlparser.AndExpr:
		return e.AndExpr(&expr, tables, context, isFunctionParamCompiler, extractAlias, applyContext)
	case *sqlparser.OrExpr:
		return e.OrExpr(expr, tables, context, isFunctionParamCompiler, extractAlias, applyContext)
	case sqlparser.OrExpr:
		return e.OrExpr(&expr, tables, context, isFunctionParamCompiler, extractAlias, applyContext)
	case *sqlparser.ComparisonExpr:
		return e.ComparisonExpr(expr, tables, context, isFunctionParamCompiler, extractAlias, applyContext)
	case sqlparser.ComparisonExpr:
		return e.ComparisonExpr(&expr, tables, context, isFunctionParamCompiler, extractAlias, applyContext)
	default:
		return nil, fmt.Errorf("not support %s in file orm/expression.getField.go", expr)
	}

}
