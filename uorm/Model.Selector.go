package uorm

import (
	"fmt"
	"strings"
)

type modelsSelectors struct {
	fields []string
	table  string
}

func (m *Model) Selector(fields ...interface{}) *modelsSelectors {
	ret := &modelsSelectors{
		fields: make([]string, 0, len(fields)),
		table:  m.table.name,
	}
	for _, field := range fields {
		switch v := field.(type) {
		case string:
			ret.fields = append(ret.fields, v)
		case Field:
			ret.fields = append(ret.fields, v.expr)
		case *Field:
			ret.fields = append(ret.fields, v.expr)
		default:
			panic(fmt.Sprintf("Selector: unsupported field type %T", field))
		}
	}
	return ret
}
func (s *modelsSelectors) String() string {
	return fmt.Sprintf("SELECT %s FROM %s", strings.Join(s.fields, ", "), s.table)
}
