package orm

import "reflect"

type dbField struct {
	Name  string
	Table string
	field reflect.StructField
}

type aliasField struct {
	UnderField interface{}
	Alias      string
}

func (f *dbField) clone() *dbField {
	return &dbField{
		Name:  f.Name,
		Table: f.Table,
		field: f.field,
	}
}
func (f *dbField) As(name string) *aliasField {
	return &aliasField{
		UnderField: f,
		Alias:      name,
	}

}
