package expr

import (
	"github.com/xwb1989/sqlparser"
)

func (e *expression) compileSQLVal(v *sqlparser.SQLVal) (*expressionCompileResult, error) {
	if v.Type == sqlparser.StrVal {
		retStr := string(v.Val)
		if e.DbDriver == DB_TYPE_MSSQL {
			retStr = "N'" + retStr + "'"
			return &expressionCompileResult{Syntax: retStr}, nil

		}
		if e.DbDriver == DB_TYPE_POSTGRES {
			retStr = "'" + "'" + retStr + "'::citext" + "'"
			return &expressionCompileResult{Syntax: retStr}, nil
		}
		return &expressionCompileResult{Syntax: retStr}, nil
	}

	return &expressionCompileResult{Syntax: "?"}, nil
	//panic(fmt.Sprintf("unsupported SQLVal type: %d, file orm/expr//expressiom.SQLVal.go", v.Type))
}
