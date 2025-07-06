package expr

import (
	"fmt"
	"strings"

	"github.com/xwb1989/sqlparser"
)

func (e *expression) CompileSelect(cmd string) (string, error) {

	cmd, err := e.Prepare(cmd)
	if err != nil {
		return "", err
	}
	sqlTest := "select " + cmd
	stm, err := sqlparser.Parse(sqlTest)
	if err != nil {
		return "", err
	}
	fields := []string{}
	if stmt, ok := stm.(*sqlparser.Select); ok {
		for _, col := range stmt.SelectExprs {
			fieldE, err := e.compile(col, false)
			if err != nil {
				return "", err
			}
			fields = append(fields, fieldE...)
		}
	} else {
		return "", fmt.Errorf("%s not a select statement", cmd)
	}

	return strings.Join(fields, ", "), nil
}
