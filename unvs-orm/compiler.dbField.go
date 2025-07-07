package orm

import "errors"

func (c *CompilerUtils) resolveDBField(context *map[string]string, f *dbField) (*resolverResult, error) {
	if f == nil {
		return nil, errors.New("dbField is nil")
	}
	if context == nil {
		return &resolverResult{
			Syntax: c.Quote(f.Table, f.Name),
			Args:   nil,
		}, nil
	}
	if alias, ok := (*context)[f.Table]; ok {
		return &resolverResult{
			Syntax:       c.Quote(alias, f.Name),
			Args:         nil,
			buildContext: context,
		}, nil
	}
	return &resolverResult{
		Syntax: c.Quote(f.Table, f.Name),
		Args:   nil,
	}, nil
}
