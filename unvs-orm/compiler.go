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
	Tables       *[]string
	hasNewTable  bool
	NewTableName string
	IsJoinExpr   bool
	NextJoin     string
}

func (r *resolverResult) GetTableAliasMap() map[string]string {
	return *r.buildContext
}

type CompilerUtils struct {
	dialect    DialectCompiler
	cacheQuote sync.Map
}

var cacheCompilerUtilsCtx sync.Map

func (c *CompilerUtils) Quote(args ...string) string {
	key := strings.Join(args, ".")
	if v, ok := c.cacheQuote.Load(key); ok {
		return v.(string)
	}

	quoteIdent := c.dialect.getQuoteIdent()
	left := string(quoteIdent[0])
	right := string(quoteIdent[1])
	ret := left + strings.Join(args, right+"."+left) + right
	c.cacheQuote.Store(key, ret)
	return ret
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
func (c *CompilerUtils) ResolveWithoutTableAlias(expr interface{}) (*resolverResult, error) {
	context := make(map[string]string)
	tables := make([]string, 0)
	return c.Resolve(&tables, &context, expr, false)
}
func (c *CompilerUtils) ResolveWithTableAlias(expr interface{}) (*resolverResult, error) {
	context := make(map[string]string)
	tables := make([]string, 0)
	r, e := c.Resolve(&tables, &context, expr, true)
	if e != nil {
		return nil, e
	}
	r.buildContext = &context
	return r, nil
}

func (c *CompilerUtils) Resolve(tables *[]string, context *map[string]string, expr interface{}, requireAlias bool) (*resolverResult, error) {
	if expr == nil {
		return &resolverResult{
			Syntax: "",
			Args:   nil,
		}, nil
	}
	typ := reflect.TypeOf(expr)
	if typ.Kind() == reflect.Slice {
		args, err := c.resolveSlice(tables, context, expr, requireAlias)
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
		return c.resolveDBField(tables, context, f, requireAlias)
	}
	if f, ok := expr.(dbField); ok {
		return c.resolveDBField(tables, context, &f, requireAlias)
	}
	if f, ok := expr.(*aliasField); ok {
		result, err := c.Resolve(tables, context, f.underField, requireAlias)
		if err != nil {
			return nil, err
		}
		return &resolverResult{
			Syntax: fmt.Sprintf("%s AS %s", result.Syntax, c.Quote(f.Alias)),
			Args:   result.Args,
		}, nil
	}

	if f, ok := expr.(aliasField); ok {
		return c.Resolve(tables, context, &f, requireAlias)
	}
	if f, ok := expr.(*fieldBinary); ok {
		return c.resolveBinaryField(tables, context, f, requireAlias)

	}

	if f, ok := expr.(fieldBinary); ok {
		return c.Resolve(tables, context, &f, requireAlias)
	}
	if f, ok := expr.(*BoolField); ok {
		return c.resolveBoolField(tables, context, f, requireAlias)

	}
	if f, ok := expr.(BoolField); ok {
		return c.resolveBoolField(tables, context, &f, requireAlias)
	}
	if f, ok := expr.(*DateTimeField); ok {
		// if f.callMethod != nil {
		// 	return c.dialect.resolve(context, f.callMethod) //<-- call method resolver no longer refers to the Field
		// }
		return c.Resolve(tables, context, f.underField, requireAlias)
	}
	if f, ok := expr.(DateTimeField); ok {

		return c.Resolve(tables, context, f.underField, requireAlias)
	}
	if f, ok := expr.(*NumberField[int]); ok {
		// if f.callMethod != nil {
		// 	return c.dialect.resolve(context, f.callMethod) //<-- call method resolver no longer refers to the Field

		// }
		return c.Resolve(tables, context, f.underField, requireAlias)

	}
	if f, ok := expr.(NumberField[int]); ok {
		return c.Resolve(tables, context, f.underField, requireAlias)
	}
	if f, ok := expr.(*TextField); ok {
		return c.Resolve(tables, context, f.underField, requireAlias)

	}
	if f, ok := expr.(TextField); ok {
		return c.Resolve(tables, context, f.underField, requireAlias)
	}
	if f, ok := expr.(*methodCall); ok {
		return c.dialect.resolve(tables, context, f, requireAlias) //<-- this method will compile with alias if any table of field found in context
	}
	if f, ok := expr.(methodCall); ok {
		return c.dialect.resolve(tables, context, &f, requireAlias)
	}
	if f, ok := expr.(*expression.MethodCall); ok {
		return c.dialect.resolve(tables, context, &methodCall{
			method:    f.Method,
			args:      f.Args,
			isFromExr: true,
		}, requireAlias)
	}
	if f, ok := expr.(expression.MethodCall); ok {
		return c.dialect.resolve(tables, context, &methodCall{
			method: f.Method,
			args:   f.Args,
		}, requireAlias)
	}
	if f, ok := expr.(*exprField); ok {
		return c.resolveExprField(tables, context, f, requireAlias)
	}
	if f, ok := expr.(exprField); ok {
		return c.resolveExprField(tables, context, &f, requireAlias)
	}
	if f, ok := expr.(*joinField); ok {
		return c.resolveJoinField(tables, context, *f, requireAlias)
	}
	if f, ok := expr.(joinField); ok {
		return c.resolveJoinField(tables, context, f, requireAlias)
	}
	ret, err := c.resolveNumberField(tables, context, expr, requireAlias)
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
func (c *CompilerUtils) ResolveBetween(tables *[]string, context *map[string]string, btf *BoolField, requireAlias bool) (*resolverResult, error) {
	if f, ok := btf.underField.(*fieldBinary); ok {
		left, err := c.Resolve(tables, context, f.left, requireAlias)
		if err != nil {
			return nil, err
		}
		right, err := c.Resolve(tables, context, f.right, requireAlias)
		if err != nil {
			return nil, err
		}
		c.addMultipleTables(tables, context, left.Tables, right.Tables)
		args := append(left.Args, right.Args...)
		right.Syntax = strings.ReplaceAll(right.Syntax, ",", " AND ")
		return &resolverResult{
			Syntax: fmt.Sprintf("%s BETWEEN %s", left.Syntax, right.Syntax),
			Args:   args,
		}, nil
	}
	return nil, fmt.Errorf("unsupported expression type: %T, file %s, line %d", btf.underField, "unvs-orm/compiler.go", 232)
}
func (c *CompilerUtils) ResolveNotBetween(tables *[]string, context *map[string]string, btf *BoolField, requireAlias bool) (*resolverResult, error) {
	if f, ok := btf.underField.(*fieldBinary); ok {
		left, err := c.Resolve(tables, context, f.left, requireAlias)
		if err != nil {
			return nil, err
		}
		right, err := c.Resolve(tables, context, f.right, requireAlias)
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
	return nil, fmt.Errorf("unsupported expression type: %T, file %s, line %d", btf, "unvs-orm/compiler.go", 251)
}
func (c *CompilerUtils) resolveSlice(tables *[]string, context *map[string]string, expr interface{}, requireAlias bool) ([]*resolverResult, error) {
	slice := reflect.ValueOf(expr)
	results := make([]*resolverResult, slice.Len())
	for i := 0; i < slice.Len(); i++ {
		result, err := c.Resolve(tables, context, slice.Index(i).Interface(), requireAlias)
		if err != nil {
			return nil, err
		}
		results[i] = result
	}
	return results, nil
}

func (c *CompilerUtils) resolveBoolField(tables *[]string, context *map[string]string, bf *BoolField, requireAlias bool) (*resolverResult, error) {
	if _, ok := bf.underField.(*joinField); ok {
		return c.resolveBoolFieldJoin(tables, context, bf, requireAlias)
	}

	return c.Resolve(tables, context, bf.underField, requireAlias)
	//return nil, fmt.Errorf("unsupported expression type: %T, file %s, line %d", bf.underField, "unvs-orm/compiler.go", 317)
}

var Compiler = CompilerUtils{}
