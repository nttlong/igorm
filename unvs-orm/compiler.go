package orm

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type resolverResult struct {
	Syntax string
	Args   []interface{}
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
func (c *CompilerUtils) Resolve(expr interface{}) (*resolverResult, error) {
	if expr == nil {
		return &resolverResult{
			Syntax: "",
			Args:   nil,
		}, nil
	}
	typ := reflect.TypeOf(expr)
	if typ.Kind() == reflect.Slice {
		args, err := c.resolveSlice(expr)
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
		return &resolverResult{
			Syntax: c.Quote(f.Table, f.Name),
			Args:   nil,
		}, nil
	}
	if f, ok := expr.(dbField); ok {
		return c.Resolve(&f)
	}
	if f, ok := expr.(*aliasField); ok {
		result, err := c.Resolve(f.Expr)
		if err != nil {
			return nil, err
		}
		return &resolverResult{
			Syntax: fmt.Sprintf("%s AS %s", result.Syntax, c.Quote(f.Alias)),
			Args:   result.Args,
		}, nil
	}

	if f, ok := expr.(aliasField); ok {
		return c.Resolve(&f)
	}
	if f, ok := expr.(*fieldBinary); ok {
		left, err := c.Resolve(f.left)
		if err != nil {
			return nil, err
		}
		right, err := c.Resolve(f.right)
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
		return c.Resolve(&f)
	}
	if f, ok := expr.(*BoolField); ok {
		return c.resolveBoolField(f)

	}
	if f, ok := expr.(BoolField); ok {
		return c.resolveBoolField(&f)
	}
	if f, ok := expr.(*DateTimeField); ok {
		if f.callMethod != nil {
			return c.dialect.resolve(f.callMethod) //<-- call method resolver no longer refers to the Field
		}
		return c.Resolve(f.dbField)
	}
	if f, ok := expr.(DateTimeField); ok {

		return c.Resolve(f.dbField)
	}
	if f, ok := expr.(*NumberField[int]); ok {
		if f.callMethod != nil {
			return c.dialect.resolve(f.callMethod) //<-- call method resolver no longer refers to the Field

		}
		return c.Resolve(f.dbField)
	}
	if f, ok := expr.(NumberField[int]); ok {
		return c.Resolve(f.dbField)
	}
	if f, ok := expr.(*TextField); ok {
		if f.callMethod != nil {
			return c.dialect.resolve(f.callMethod) //<-- call method resolver no longer refers to the Field
		}
		return c.Resolve(f.dbField)
	}
	if f, ok := expr.(TextField); ok {
		return c.Resolve(f.dbField)
	}
	if f, ok := expr.(*methodCall); ok {
		return c.dialect.resolve(f)
	}
	if f, ok := expr.(methodCall); ok {
		return c.dialect.resolve(&f)
	}
	ret, err := c.resolveNumberField(expr)
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

	panic(fmt.Errorf("unsupported expression type: %T, file %s, line %d", expr, "compiler.go", 179))
}
func (c *CompilerUtils) ResolveBetween(f *BoolField) (*resolverResult, error) {
	left, err := c.Resolve(f.left)
	if err != nil {
		return nil, err
	}
	right, err := c.Resolve(f.right)
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
func (c *CompilerUtils) ResolveNotBetween(f *BoolField) (*resolverResult, error) {
	left, err := c.Resolve(f.left)
	if err != nil {
		return nil, err
	}
	right, err := c.Resolve(f.right)
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
func (c *CompilerUtils) resolveSlice(expr interface{}) ([]*resolverResult, error) {
	slice := reflect.ValueOf(expr)
	results := make([]*resolverResult, slice.Len())
	for i := 0; i < slice.Len(); i++ {
		result, err := c.Resolve(slice.Index(i).Interface())
		if err != nil {
			return nil, err
		}
		results[i] = result
	}
	return results, nil
}

func (c *CompilerUtils) resolveNumberField(expr interface{}) (*resolverResult, error) {
	if f, ok := expr.(*NumberField[int64]); ok {
		return c.Resolve(f.dbField)
	}
	if f, ok := expr.(NumberField[int64]); ok {
		return c.Resolve(&f)
	}
	if f, ok := expr.(*NumberField[float64]); ok {
		return c.Resolve(f.dbField)
	}
	if f, ok := expr.(NumberField[float64]); ok {
		return c.Resolve(&f)
	}
	if f, ok := expr.(*NumberField[int64]); ok {
		return c.Resolve(f.dbField)
	}
	if f, ok := expr.(NumberField[int64]); ok {
		return c.Resolve(&f)
	}
	if f, ok := expr.(*NumberField[float64]); ok {
		return c.Resolve(f.dbField)
	}
	if f, ok := expr.(NumberField[float64]); ok {
		return c.Resolve(&f)
	}
	return nil, nil
}
func (c *CompilerUtils) resolveConstant(expr interface{}) (*resolverResult, error) {
	switch expr.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, string:
		return &resolverResult{
			Syntax: "?",
			Args:   []interface{}{expr},
		}, nil
	}
	return nil, nil
}
func (c *CompilerUtils) resolveBoolField(f *BoolField) (*resolverResult, error) {
	if f.op == "BETWEEN" {
		return c.ResolveBetween(f)

	}
	if f.op == "NOT BETWEEN" {
		return c.ResolveNotBetween(f)

	}

	left, err := c.Resolve(f.left)
	if err != nil {
		return nil, err
	}
	right, err := c.Resolve(f.right)
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
