package orm

import (
	"strings"
	"sync"
)

type JoinCompilerUtils struct {
	dialect DialectCompiler
}

var cacheJoinCompilerCtx sync.Map

func (c *JoinCompilerUtils) Ctx(dialect DialectCompiler) *JoinCompilerUtils {
	if dialect == nil {
		panic("dialect is nil")
	}

	key := dialect.driverName()

	if v, ok := cacheJoinCompilerCtx.Load(key); ok {
		return v.(*JoinCompilerUtils)
	}
	ret := &JoinCompilerUtils{dialect: dialect}
	dialect.setJoinCompiler(ret)
	cacheJoinCompilerCtx.Store(key, ret)
	return ret
}
func (c *JoinCompilerUtils) Resolve(expr *JoinExpr) (*resolverResult, error) {
	cmp := Compiler.Ctx(c.dialect) //<-- get compiler for dialect
	retSyntax := []string{}
	args := []interface{}{}
	// stack := []*JoinExpr{}
	for node := expr; node != nil; node = node.previous {
		if node.on != nil {
			r, err := cmp.Resolve(&node.aliasMap, node.on) //<-- resolve on condition
			if err != nil {
				return nil, err
			}
			syntax := cmp.Quote(node.baseTable) + " AS " + cmp.Quote(node.aliasMap[node.baseTable]) + " ON " + r.Syntax
			if node.previous != nil {
				if node.joinType == "" {
					panic("joinType is empty")
				}
				syntax = node.previous.joinType + " JOIN " + syntax
			}
			retSyntax = append([]string{syntax}, retSyntax...)
			args = append(r.Args, args...)

		} else {
			syntax := cmp.Quote(node.baseTable) + " AS " + cmp.Quote(node.aliasMap[node.baseTable])
			retSyntax = append([]string{syntax}, retSyntax...)
		}

	}

	return &resolverResult{
		Syntax:      strings.Join(retSyntax, " "),
		Args:        args,
		AliasSource: expr.aliasMap,
	}, nil

}

func (c *JoinCompilerUtils) ResolveBoolFieldAsJoin(expr *BoolField) (*resolverResult, error) {
	cmp := Compiler.Ctx(c.dialect) //<-- get compiler for dialect
	if len(expr.tables) == 0 {
		expr = expr.doJoin()
		cmpRes, err := cmp.Resolve(&expr.alias, expr.left)
		if err != nil {
			return nil, err
		}
		right := cmp.Quote(expr.tables[1]) + " AS " + cmp.Quote(expr.alias[expr.tables[1]])
		left := cmp.Quote(expr.tables[0]) + " AS " + cmp.Quote(expr.alias[expr.tables[0]])
		cmpRes.Syntax = left + " " + expr.joinType + " JOIN " + right + " ON " + cmpRes.Syntax
		return &resolverResult{
			Syntax:      cmpRes.Syntax,
			Args:        cmpRes.Args,
			AliasSource: expr.alias,
		}, nil
	}
	cmpRes, err := cmp.Resolve(&expr.alias, expr)
	if err != nil {
		return nil, err
	}
	return &resolverResult{
		Syntax:      cmpRes.Syntax,
		Args:        cmpRes.Args,
		AliasSource: expr.alias,
	}, nil

}

// stack := []*JoinExpr{}

var JoinCompiler = JoinCompilerUtils{}
