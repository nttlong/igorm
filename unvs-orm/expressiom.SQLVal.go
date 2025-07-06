package orm

import (
	"fmt"

	"github.com/xwb1989/sqlparser"
)

func (e *expression) compileSQLVal(v *sqlparser.SQLVal) ([]string, error) {
	if v.Type == sqlparser.StrVal {
		return []string{string(v.Val)}, nil
	}
	if v.Type == sqlparser.IntVal {
		return []string{string(v.Val)}, nil
	}

	panic(fmt.Sprintf("unsupported SQLVal type: %d, file orm/expressiom.SQLVal.go", v.Type))
}
