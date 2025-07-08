package orm

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

type JoinCompilerUtils struct {
	dialect             DialectCompiler
	cacheFromExprString sync.Map
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
func (c *JoinCompilerUtils) Compile(expr *BoolField) (*resolverResult, error) {
	if je, ok := expr.underField.(*JoinExpr); ok {
		return c.Resolve(je)
	} else {
		return nil, fmt.Errorf("unsupported expression type: %T", expr.underField)
	}
}
func (c *JoinCompilerUtils) Resolve(expr *JoinExpr) (*resolverResult, error) {

	if expr.joinExprText != nil {
		join, err := c.fromExprString(expr)

		if err != nil {
			return nil, err
		}

		ctx := *join.buildContext
		tables := *join.Tables
		strLeft := c.dialect.getCompiler().Quote(tables[0]) + " AS " + c.dialect.getCompiler().Quote(ctx[tables[0]])
		strRight := c.dialect.getCompiler().Quote(tables[1]) + " AS " + c.dialect.getCompiler().Quote(ctx[tables[1]])
		syntax := strLeft + " " + expr.joinType + " JOIN " + strRight + " ON " + join.Syntax
		return &resolverResult{
			Syntax:       syntax,
			Args:         join.Args,
			buildContext: join.buildContext,
		}, nil

	}
	cmp := Compiler.Ctx(c.dialect) //<-- get compiler for dialect
	retSyntax := []string{}
	args := []interface{}{}
	// stack := []*JoinExpr{}
	context := map[string]string{}
	tables := []string{}
	for node := expr; node != nil; node = node.previous {
		if node.on != nil {

			r, err := cmp.Resolve(&tables, &context, node.on, true) //<-- resolve on condition is alway required alias mapping
			if err != nil {
				return nil, err
			}

			tableName := Utils.GetTableNameFromVirtualName(node.baseTable)
			syntax := cmp.Quote(tableName) + " AS " + cmp.Quote(context[node.baseTable]) + " ON " + r.Syntax
			if node.previous != nil {
				if node.joinType == "" {
					panic("joinType is empty")
				}
				syntax = node.previous.joinType + " JOIN " + syntax
			}
			retSyntax = append([]string{syntax}, retSyntax...)
			args = append(r.Args, args...)

		} else {
			table := Utils.GetTableNameFromVirtualName(node.baseTable)

			syntax := cmp.Quote(table) + " AS " + cmp.Quote(context[node.baseTable])
			retSyntax = append([]string{syntax}, retSyntax...)
		}

	}
	ret := strings.Join(retSyntax, " ")
	return &resolverResult{
		Syntax:       ret,
		Args:         args,
		buildContext: &context,
	}, nil

}
func (c *JoinCompilerUtils) resolveJoinField(tables *[]string, context *map[string]string, expr joinField) (*resolverResult, error) {
	cmp := Compiler.Ctx(c.dialect) //<-- get compiler for dialect

	cmpRes, err := cmp.Resolve(tables, context, expr, true)
	if err != nil {
		return nil, err
	}
	return &resolverResult{
		Syntax:       cmpRes.Syntax,
		Args:         cmpRes.Args,
		buildContext: cmpRes.buildContext,
	}, nil

}
func (c *JoinCompilerUtils) resoleFieldBinary(tables *[]string, context *map[string]string, bF *BoolField, expr fieldBinary) (*resolverResult, error) {
	cmp := Compiler.Ctx(c.dialect) //<-- get compiler for dialect
	cmpRes, err := cmp.Resolve(tables, context, expr, true)
	if len(*tables) == 1 {
		leftContext := map[string]string{
			(*tables)[0]: "T" + strconv.Itoa(len(*tables)),
		}
		cmpResLeft, err := cmp.Resolve(tables, &leftContext, expr.left, true)
		if err != nil {
			return nil, err
		}
		fmt.Println(cmpResLeft.Syntax)
		rightContext := map[string]string{
			(*tables)[0]: "T" + strconv.Itoa(len(*tables)+1),
		}
		cmpResRight, err := cmp.Resolve(tables, &rightContext, expr.right, true)
		if err != nil {
			return nil, err
		}
		fmt.Println(cmpResRight.Syntax)
		syntax := cmpResLeft.Syntax + " " + expr.op + " " + cmpResRight.Syntax
		cmpRes.Syntax = syntax
		cmpRes.Args = append(cmpResLeft.Args, cmpResRight.Args...)
	}
	if err != nil {
		return nil, err
	}
	last2Table := (*tables)[len(*tables)-2:]
	leftTable := (last2Table)[0]
	rightTable := (last2Table)[1]
	leftAlias := (*context)[leftTable]
	rightAlias := (*context)[rightTable]
	left := cmp.Quote(leftTable) + " AS " + cmp.Quote(leftAlias)
	right := cmp.Quote(rightTable) + " AS " + cmp.Quote(rightAlias)
	syntax := left + " INNER JOIN " + right + " ON " + cmpRes.Syntax
	cmpRes.Syntax = syntax

	return cmpRes, nil
}
func (c *JoinCompilerUtils) ResolveBoolFieldAsJoin(tables *[]string, context *map[string]string, bF *BoolField) (*resolverResult, error) {
	if tables == nil {
		tables = &[]string{}
	}
	if context == nil {
		context = &map[string]string{}
	}
	if expr, ok := bF.underField.(*joinField); ok {

		return c.resolveJoinField(tables, context, *expr)
	}
	if expr, ok := bF.underField.(joinField); ok {
		return c.resolveJoinField(tables, context, expr)
	}
	if expr, ok := bF.underField.(fieldBinary); ok {
		return c.resoleFieldBinary(tables, context, bF, expr)
	}
	if expr, ok := bF.underField.(*fieldBinary); ok {
		return c.resoleFieldBinary(tables, context, bF, *expr)
	}
	panic(fmt.Errorf("unsupported expression type: %T, file %s, line %d", bF.underField, "unvs-orm/SqlBuildrJoin.go", 127))

}

// stack := []*JoinExpr{}

var JoinCompiler = JoinCompilerUtils{}
