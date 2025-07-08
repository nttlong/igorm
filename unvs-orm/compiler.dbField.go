package orm

import (
	"errors"
	"strconv"
)

func (c *CompilerUtils) addTables(tables *[]string, context *map[string]string, tablesAdd ...string) bool {
	hasNew := false
	for _, t := range tablesAdd {
		if _, ok := (*context)[t]; !ok {
			(*context)[t] = "T" + strconv.Itoa(len(*context)+1)
			*tables = append(*tables, t)
			hasNew = true
		}
	}
	return hasNew
}
func (c *CompilerUtils) addMultipleTables(tables *[]string, context *map[string]string, tablesAdd ...*[]string) {
	for _, t := range tablesAdd {
		if t != nil {
			c.addTables(tables, context, (*t)...)
		}
	}
}
func (c *CompilerUtils) resolveDBField(tables *[]string, context *map[string]string, f *dbField, requireAlias bool) (*resolverResult, error) {
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
	hasNew := c.addTables(tables, context, f.Table)
	if requireAlias {
		return &resolverResult{
			Syntax:       c.Quote((*context)[f.Table], f.Name),
			Tables:       tables,
			Args:         nil,
			buildContext: context,
			hasNewTable:  hasNew,
			NewTableName: f.Table,
		}, nil
	} else {
		return &resolverResult{
			Syntax:       c.Quote(f.Table, f.Name),
			Tables:       tables,
			Args:         nil,
			buildContext: context,
		}, nil
	}

}
