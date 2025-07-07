package orm

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	expression "unvs-orm/expr"
)

type resolverResult struct {
	Syntax       string
	Args         []interface{}
	buildContext *map[string]string
	Tables       []string
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
func (c *CompilerUtils) Resolve(context *map[string]string, expr interface{}) (*resolverResult, error) {
	if expr == nil {
		return &resolverResult{
			Syntax: "",
			Args:   nil,
		}, nil
	}
	typ := reflect.TypeOf(expr)
	if typ.Kind() == reflect.Slice {
		args, err := c.resolveSlice(context, expr)
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
		return c.resolveDBField(context, f)
	}
	if f, ok := expr.(dbField); ok {
		return c.resolveDBField(context, &f)
	}
	if f, ok := expr.(*aliasField); ok {
		result, err := c.Resolve(context, f.UnderField)
		if err != nil {
			return nil, err
		}
		return &resolverResult{
			Syntax: fmt.Sprintf("%s AS %s", result.Syntax, c.Quote(f.Alias)),
			Args:   result.Args,
		}, nil
	}

	if f, ok := expr.(aliasField); ok {
		return c.Resolve(context, &f)
	}
	if f, ok := expr.(*fieldBinary); ok {
		return c.resolveBinaryField(context, f)

	}

	if f, ok := expr.(fieldBinary); ok {
		return c.Resolve(context, &f)
	}
	if f, ok := expr.(*BoolField); ok {
		return c.resolveBoolField(context, f)

	}
	if f, ok := expr.(BoolField); ok {
		return c.resolveBoolField(context, &f)
	}
	if f, ok := expr.(*DateTimeField); ok {
		// if f.callMethod != nil {
		// 	return c.dialect.resolve(context, f.callMethod) //<-- call method resolver no longer refers to the Field
		// }
		return c.Resolve(context, f.UnderField)
	}
	if f, ok := expr.(DateTimeField); ok {

		return c.Resolve(context, f.UnderField)
	}
	if f, ok := expr.(*NumberField[int]); ok {
		// if f.callMethod != nil {
		// 	return c.dialect.resolve(context, f.callMethod) //<-- call method resolver no longer refers to the Field

		// }
		return c.Resolve(context, f.UnderField)

	}
	if f, ok := expr.(NumberField[int]); ok {
		return c.Resolve(context, f.UnderField)
	}
	if f, ok := expr.(*TextField); ok {
		return c.Resolve(context, f.UnderField)

	}
	if f, ok := expr.(TextField); ok {
		return c.Resolve(context, f.UnderField)
	}
	if f, ok := expr.(*methodCall); ok {
		return c.dialect.resolve(context, f)
	}
	if f, ok := expr.(methodCall); ok {
		return c.dialect.resolve(context, &f)
	}
	if f, ok := expr.(*expression.MethodCall); ok {
		return c.dialect.resolve(context, &methodCall{
			method: f.Method,
			args:   f.Args,
		})
	}
	if f, ok := expr.(expression.MethodCall); ok {
		return c.dialect.resolve(context, &methodCall{
			method: f.Method,
			args:   f.Args,
		})
	}
	if f, ok := expr.(*exprField); ok {
		return c.resolveExprField(context, f)
	}
	if f, ok := expr.(exprField); ok {
		return c.resolveExprField(context, &f)
	}
	if f, ok := expr.(*joinField); ok {
		return c.resolveJoinField(context, *f)
	}
	if f, ok := expr.(joinField); ok {
		return c.resolveJoinField(context, f)
	}
	ret, err := c.resolveNumberField(context, expr)
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
func (c *CompilerUtils) ResolveBetween(context *map[string]string, btf *BoolField) (*resolverResult, error) {
	if f, ok := btf.UnderField.(*fieldBinary); ok {
		left, err := c.Resolve(context, f.left)
		if err != nil {
			return nil, err
		}
		right, err := c.Resolve(context, f.right)
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
	return nil, fmt.Errorf("unsupported expression type: %T, file %s, line %d", btf.UnderField)
}
func (c *CompilerUtils) ResolveNotBetween(context *map[string]string, btf *BoolField) (*resolverResult, error) {
	if f, ok := btf.UnderField.(*fieldBinary); ok {
		left, err := c.Resolve(context, f.left)
		if err != nil {
			return nil, err
		}
		right, err := c.Resolve(context, f.right)
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
	return nil, fmt.Errorf("unsupported expression type: %T, file %s, line %d", btf)
}
func (c *CompilerUtils) resolveSlice(context *map[string]string, expr interface{}) ([]*resolverResult, error) {
	slice := reflect.ValueOf(expr)
	results := make([]*resolverResult, slice.Len())
	for i := 0; i < slice.Len(); i++ {
		result, err := c.Resolve(context, slice.Index(i).Interface())
		if err != nil {
			return nil, err
		}
		results[i] = result
	}
	return results, nil
}

func (c *CompilerUtils) resolveBoolField(context *map[string]string, bf *BoolField) (*resolverResult, error) {
	if _, ok := bf.UnderField.(*joinField); ok {
		return c.resolveBoolFieldJoin(context, bf)
	}

	return c.Resolve(context, bf.UnderField)
	//return nil, fmt.Errorf("unsupported expression type: %T, file %s, line %d", bf.UnderField, "unvs-orm/compiler.go", 317)
}

var Compiler = CompilerUtils{}
