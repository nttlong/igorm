package orm

import "fmt"

func (c *CompilerUtils) resolveBoolFieldJoin(tables *[]string, context *map[string]string, bf *BoolField, requireAlias bool) (*resolverResult, error) {
	if f, ok := bf.underField.(*joinField); ok {
		left, err := c.Resolve(tables, context, f.left, requireAlias)
		if err != nil {
			return nil, err
		}
		right, err := c.Resolve(tables, context, f.right, requireAlias)
		if err != nil {
			return nil, err
		}
		args := append(left.Args, right.Args...)
		rightSource := c.Quote(f.tables[0]) + " AS " + c.Quote(f.alias[f.tables[0]])
		syntax := left.Syntax + " " + f.joinType + " " + rightSource + " ON " + right.Syntax
		return &resolverResult{
			Syntax: syntax, //fmt.Sprintf("%s %s %s ON %s", left.Syntax, f.op, right.Syntax, f.rawText),
			Args:   args,
		}, nil
	}
	return nil, fmt.Errorf("unsupported expression type: %T, file %s, line %d", bf.underField, "unvs-orm/compiler.go", 23)
}
