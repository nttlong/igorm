package expr

import (
	"fmt"
	"strings"
	"sync"

	"github.com/xwb1989/sqlparser"
)

type compileResult struct {
	Syntax  string
	Args    []interface{}
	Context *map[string]string
}

//	func (e *expression) CompileSelect(cmd string) (*compileResult, error) {
//		return e.CompileSelectFull(nil, cmd, false)
//	}
func (e *expression) Compile(tables *[]string, context *map[string]string, cmd string, requireAlias bool) (*compileResult, error) {

	cmd = e.Prepare(cmd)

	sqlTest := "select " + cmd
	stm, err := sqlparser.Parse(sqlTest)
	if err != nil {
		return nil, err
	}
	fields := []string{}

	if stmt, ok := stm.(*sqlparser.Select); ok {
		for _, col := range stmt.SelectExprs {

			fieldE, err := e.compile(col, tables, context, false, true)
			if err != nil {
				return nil, err
			}
			fields = append(fields, fieldE.Syntax)
		}
	} else {
		return nil, fmt.Errorf("%s not a select statement", cmd)
	}

	syntax := strings.Join(fields, ", ")
	return &compileResult{
		Syntax:  syntax,
		Args:    []interface{}{},
		Context: context,
	}, nil
}

var cacheNewExpressionCompiler sync.Map

func NewExpressionCompiler(driver string) *expression {
	if v, ok := cacheNewExpressionCompiler.Load(driver); ok {
		return v.(*expression)
	}
	e := &expression{
		keywords:    nil,
		specialChar: []byte{'.', ' ', '\t', '\n', '\r', ',', ';', '(', ')', '<', '>', '=', '+', '-', '*', '/', '%', '&', '|', '^', '!', '?'},
		DbDriver:    DB_TYPE_UNKNOWN.FromString(driver),
	}
	cacheNewExpressionCompiler.Store(driver, e)
	return e
}
