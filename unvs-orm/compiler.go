package orm

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type resolverResult struct {
	Syntax      string
	Args        []interface{}
	AliasSource map[string]string
}
type CompilerUtils struct {
	dialect DialectCompiler
}

var cacheCompilerUtilsCtx sync.Map

func (c *CompilerUtils) Quote(args ...string) string {
	quoteIdent := c.dialect.getQuoteIdent()
	left := string(quoteIdent[0])
	right := string(quoteIdent[1])
	return left + strings.Join(args, right+"."+left) + right

}
func (c *CompilerUtils) Ctx(dialect DialectCompiler) *CompilerUtils {
	if dialect == nil {
		panic("dialect is nil")
	}

	key := dialect.driverName()

	if v, ok := cacheCompilerUtilsCtx.Load(key); ok {
		return v.(*CompilerUtils)
	}
	ret := &CompilerUtils{dialect: dialect}
	dialect.setCompiler(ret)
	cacheCompilerUtilsCtx.Store(key, ret)
	return ret
}
func (c *CompilerUtils) Resolve(aliasSource *map[string]string, expr interface{}) (*resolverResult, error) {
	if expr == nil {
		return &resolverResult{
			Syntax: "",
			Args:   nil,
		}, nil
	}
	typ := reflect.TypeOf(expr)
	if typ.Kind() == reflect.Slice {
		args, err := c.resolveSlice(aliasSource, expr)
		if err != nil {
			return nil, err
		}
		ret := &resolverResult{
			Syntax: "",
			Args:   make([]interface{}, 0),
		}
		for _, arg := range args {
			if arg.Syntax == "" {
				ret.Args = append(ret.Args, arg.Args...)
			} else {
				ret.Syntax += arg.Syntax + ","
				ret.Args = append(ret.Args, arg.Args...)
			}
		}
		ret.Syntax = strings.TrimSuffix(ret.Syntax, ",")
		return ret, nil

	}
	if f, ok := expr.(*dbField); ok {
		return c.resolveDBField(aliasSource, f)
	}
	if f, ok := expr.(dbField); ok {
		return c.resolveDBField(aliasSource, &f)
	}
	if f, ok := expr.(*aliasField); ok {
		result, err := c.Resolve(aliasSource, f.Expr)
		if err != nil {
			return nil, err
		}
		return &resolverResult{
			Syntax: fmt.Sprintf("%s AS %s", result.Syntax, c.Quote(f.Alias)),
			Args:   result.Args,
		}, nil
	}

	if f, ok := expr.(aliasField); ok {
		return c.Resolve(aliasSource, &f)
	}
	if f, ok := expr.(*fieldBinary); ok {
		left, err := c.Resolve(aliasSource, f.left)
		if err != nil {
			return nil, err
		}
		right, err := c.Resolve(aliasSource, f.right)
		if err != nil {
			return nil, err
		}
		args := append(left.Args, right.Args...)
		if len(args) == 0 {
			return &resolverResult{
				Syntax: fmt.Sprintf("%s %s", left.Syntax, f.op),
				Args:   args,
			}, nil
		} else if len(args) == 1 {
			return &resolverResult{
				Syntax: fmt.Sprintf("%s %s %s", left.Syntax, f.op, right.Syntax),
				Args:   args,
			}, nil
		} else {
			return &resolverResult{
				Syntax: fmt.Sprintf("%s %s (%s)", left.Syntax, f.op, right.Syntax),
				Args:   args,
			}, nil
		}

	}

	if f, ok := expr.(fieldBinary); ok {
		return c.Resolve(aliasSource, &f)
	}
	if f, ok := expr.(*BoolField); ok {
		return c.resolveBoolField(aliasSource, f)

	}
	if f, ok := expr.(BoolField); ok {
		return c.resolveBoolField(aliasSource, &f)
	}
	if f, ok := expr.(*DateTimeField); ok {
		if f.callMethod != nil {
			return c.dialect.resolve(aliasSource, f.callMethod) //<-- call method resolver no longer refers to the Field
		}
		return c.Resolve(aliasSource, f.dbField)
	}
	if f, ok := expr.(DateTimeField); ok {

		return c.Resolve(aliasSource, f.dbField)
	}
	if f, ok := expr.(*NumberField[int]); ok {
		if f.callMethod != nil {
			return c.dialect.resolve(aliasSource, f.callMethod) //<-- call method resolver no longer refers to the Field

		}
		return c.Resolve(aliasSource, f.dbField)
	}
	if f, ok := expr.(NumberField[int]); ok {
		return c.Resolve(aliasSource, f.dbField)
	}
	if f, ok := expr.(*TextField); ok {
		if f.callMethod != nil {
			return c.dialect.resolve(aliasSource, f.callMethod) //<-- call method resolver no longer refers to the Field
		}
		return c.Resolve(aliasSource, f.dbField)
	}
	if f, ok := expr.(TextField); ok {
		return c.Resolve(aliasSource, f.dbField)
	}
	if f, ok := expr.(*methodCall); ok {
		return c.dialect.resolve(aliasSource, f)
	}
	if f, ok := expr.(methodCall); ok {
		return c.dialect.resolve(aliasSource, &f)
	}
	ret, err := c.resolveNumberField(aliasSource, expr)
	if err != nil {
		return nil, err
	}
	if ret == nil {
		ret, err := c.resolveConstant(expr)
		if err != nil {
			return nil, err
		}
		if ret != nil {
			return ret, nil
		}

	}
	if ret != nil {
		return ret, nil
	}
	// endregion
	//*unvs-orm.NumberField[int8]

	panic(fmt.Errorf("unsupported expression type: %T, file %s, line %d", expr, "unvs-orm/compiler.go", 187))
}
func (c *CompilerUtils) ResolveBetween(aliasSource *map[string]string, f *BoolField) (*resolverResult, error) {
	left, err := c.Resolve(aliasSource, f.left)
	if err != nil {
		return nil, err
	}
	right, err := c.Resolve(aliasSource, f.right)
	if err != nil {
		return nil, err
	}
	args := append(left.Args, right.Args...)
	right.Syntax = strings.ReplaceAll(right.Syntax, ",", " AND ")
	return &resolverResult{
		Syntax: fmt.Sprintf("%s BETWEEN %s", left.Syntax, right.Syntax),
		Args:   args,
	}, nil
}
func (c *CompilerUtils) ResolveNotBetween(aliasSource *map[string]string, f *BoolField) (*resolverResult, error) {
	left, err := c.Resolve(aliasSource, f.left)
	if err != nil {
		return nil, err
	}
	right, err := c.Resolve(aliasSource, f.right)
	if err != nil {
		return nil, err
	}
	args := append(left.Args, right.Args...)
	right.Syntax = strings.ReplaceAll(right.Syntax, ",", " AND ")
	return &resolverResult{
		Syntax: fmt.Sprintf("%s NOT BETWEEN %s", left.Syntax, right.Syntax),
		Args:   args,
	}, nil
}
func (c *CompilerUtils) resolveSlice(aliasSource *map[string]string, expr interface{}) ([]*resolverResult, error) {
	slice := reflect.ValueOf(expr)
	results := make([]*resolverResult, slice.Len())
	for i := 0; i < slice.Len(); i++ {
		result, err := c.Resolve(aliasSource, slice.Index(i).Interface())
		if err != nil {
			return nil, err
		}
		results[i] = result
	}
	return results, nil
}

func (c *CompilerUtils) resolveBoolField(aliasSource *map[string]string, f *BoolField) (*resolverResult, error) {
	if strings.HasSuffix(f.op, " JOIN") {
		return c.resolveBoolFieldJoin(aliasSource, f)
	}
	if f.op == "BETWEEN" {
		return c.ResolveBetween(aliasSource, f)

	}
	if f.op == "NOT BETWEEN" {
		return c.ResolveNotBetween(aliasSource, f)

	}
	var left *resolverResult
	if f.left != nil && f.right != nil {
		_left, err := c.Resolve(aliasSource, f.left)
		if err != nil {
			return nil, err
		}
		left = _left
	}
	if left == nil {
		_left, err := c.Resolve(aliasSource, f.dbField)
		if err != nil {
			return nil, err
		}
		left = _left
	}
	right, err := c.Resolve(aliasSource, f.right)
	if err != nil {
		return nil, err
	}

	args := append(left.Args, right.Args...)
	if f.op == "IN" || f.op == "NOT IN" {
		return &resolverResult{
			Syntax: fmt.Sprintf("%s %s %s", left.Syntax, f.op, "("+right.Syntax+")"),
			Args:   args,
		}, nil
	}
	if right.Syntax != "" {
		return &resolverResult{
			Syntax: fmt.Sprintf("%s %s %s", left.Syntax, f.op, right.Syntax),
			Args:   args,
		}, nil
	} else {
		if f.op == "IS NULL" || f.op == "IS NOT NULL" {
			return &resolverResult{
				Syntax: fmt.Sprintf("%s %s", left.Syntax, f.op),
				Args:   []interface{}{},
			}, nil
		}
		return &resolverResult{
			Syntax: fmt.Sprintf("%s %s", left.Syntax, f.op),
			Args:   args,
		}, nil
	}
}

var Compiler = CompilerUtils{}
