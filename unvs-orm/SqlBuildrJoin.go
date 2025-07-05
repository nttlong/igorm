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

// func (c *JoinCompilerUtils) Resolve(expr *JoinExpr) (*resolverResult, error) {
// 	cmd := Compiler.Ctx(c.dialect)

// 	// Dịch điều kiện ON
// 	onRes, err := cmd.Resolve(&expr.aliasSource, expr.on)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Bảng được JOIN (phải là bảng bên phải)
// 	// Giả sử expr.tables luôn là [leftTable, rightTable]
// 	if len(expr.tables) != 2 {
// 		return nil, errors.New("JoinExpr expects exactly 2 tables")
// 	}

// 	// Bảng bên phải là bảng được JOIN vào FROM gốc
// 	rightTable := expr.tables[1]
// 	rightAlias := expr.aliasSource[rightTable]

// 	// Sinh JOIN clause
// 	tblJoin := cmd.Quote(rightTable) + " AS " + cmd.Quote(rightAlias)

// 	ret := &resolverResult{
// 		Syntax:      expr.joinType + " JOIN " + tblJoin + " ON " + onRes.Syntax,
// 		Args:        onRes.Args,
// 		AliasSource: expr.aliasSource,
// 	}
// 	return ret, nil
// }

var JoinCompiler = JoinCompilerUtils{}
