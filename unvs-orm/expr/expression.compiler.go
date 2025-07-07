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
	Context *ResolveContext
}

func (e *expression) CompileSelect(cmd string) (*compileResult, error) {

	cmd, err := e.Prepare(cmd)
	if err != nil {
		return nil, err
	}
	sqlTest := "select " + cmd
	stm, err := sqlparser.Parse(sqlTest)
	if err != nil {
		return nil, err
	}
	fields := []string{}
	context := &ResolveContext{
		Tables: []string{},
		Map:    map[string]string{},
	}
	if stmt, ok := stm.(*sqlparser.Select); ok {
		for _, col := range stmt.SelectExprs {

			fieldE, err := e.compile(col, context, false)
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
		keywords: []string{
			"order",
			"select",
			"from",
			"where",
			"and",
			"or",
			"not",
			"in",
			"like",
			"is",
			"null",
			"between",
			"exists",
			"case",
			"when",
			"then",
			"else",
			"end",
			"as",
			"distinct",
			"count",
			"sum",
			"avg",
			"max",
			"min",
			"abs",
			"ceil",
			"floor",
			"round",
			"length",
			"substring",
			"trim",
			"lower",
			"upper",
			"date",
			"time",
			"datetime",
			"year",
			"month",
			"day",
			"hour",
			"minute",
			"second",
			"now",
			"current_date",
			"current_time",
			"current_timestamp",
			"group",
			"order",
			"limit",
		},
		specialChar: []byte{'.', ' ', '\t', '\n', '\r', ',', ';', '(', ')', '<', '>', '=', '+', '-', '*', '/', '%', '&', '|', '^', '!', '?'},
		DbDriver:    DB_TYPE_UNKNOWN.FromString(driver),
	}
	cacheNewExpressionCompiler.Store(driver, e)
	return e
}
