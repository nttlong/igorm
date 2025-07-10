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
	Context map[string]string
	Tables  []string
}

/*
@extractAlias: extract alias from select statement and add to context return syntax with new alias of table

@applyContext:  lookup context for alias and replace with table name
*/
func (e *expression) Compile(sourceCache string, tables *[]string, context *map[string]string, cmd string, extractAlias bool, applyContext bool) (*compileResult, error) {
	//check cache
	if e == nil {
		return nil, fmt.Errorf("expression is nil")
	}

	key := sourceCache + "://" + cmd

	if !extractAlias && sourceCache != "" { //<-- only cache when extractAlias is false
		key = cmd + strings.Join(*tables, ",")
		if v, ok := e.cacheCompile.Load(key); ok {
			// ret := v.(compileResult)
			return v.(*compileResult), nil

		}
	}

	cmd = ExprPreProcessText.QuoteExpression(cmd)

	sqlTest := "select " + cmd
	stm, err := sqlparser.Parse(sqlTest)
	if err != nil {
		return nil, err
	}
	fields := []string{}

	if stmt, ok := stm.(*sqlparser.Select); ok {
		for _, col := range stmt.SelectExprs {

			fieldE, err := e.compile(col, tables, context, false, extractAlias, applyContext)
			if err != nil {
				return nil, err
			}
			fields = append(fields, fieldE.Syntax)
		}
	} else {
		return nil, fmt.Errorf("%s not a select statement", cmd)
	}

	syntax := strings.Join(fields, ", ")

	if !extractAlias && sourceCache != "" { //<-- only cache when extractAlias is false
		cacheContext := make(map[string]string)
		for k, v := range *context {
			cacheContext[k] = v
		}
		cacheTable := make([]string, len(*tables))
		copy(cacheTable, *tables)
		cacheResult := &compileResult{
			Syntax:  syntax,
			Args:    []interface{}{},
			Context: cacheContext,
			Tables:  cacheTable,
		}
		e.cacheCompile.Store(key, cacheResult)
	}
	return &compileResult{
		Syntax: syntax,
		Args:   []interface{}{},
		// Context: context,
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
