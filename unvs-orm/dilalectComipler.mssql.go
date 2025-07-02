package orm

import (
	"strconv"
	"strings"
)

type mssqlDialect struct {
	compiler *CompilerUtils
}

func (d *mssqlDialect) getQuoteIdent() string {
	return "[]"
}
func (d *mssqlDialect) getParam(index int) string {
	return "@p" + strconv.Itoa(index)

}
func (d *mssqlDialect) setCompiler(compiler *CompilerUtils) {
	d.compiler = compiler
}
func (d *mssqlDialect) getCompiler() *CompilerUtils {
	return d.compiler
}
func (d *mssqlDialect) resolve(caller *methodCall) (*resolverResult, error) {

	strArgs := make([]string, 0)
	retArgs := make([]interface{}, 0)
	if caller.dbField != nil {
		field, err := d.compiler.Resolve(caller.dbField)
		if err != nil {
			return nil, err
		}
		strArgs = append(strArgs, field.Syntax)
		retArgs = append(retArgs, field.Args...)
	}
	for _, arg := range caller.args {
		rs, err := d.compiler.Resolve(arg)
		if err != nil {
			return nil, err
		}
		strArgs = append(strArgs, rs.Syntax)
		retArgs = append(retArgs, rs.Args...)
	}

	return &resolverResult{
		Syntax: caller.method + "(" + strings.Join(strArgs, ",") + ")",
		Args:   retArgs,
	}, nil
}
func (d *mssqlDialect) driverName() string {
	return "mssql"
}

var MssqlDialect = mssqlDialect{}

func NewMssqlDialect() DialectCompiler {
	return &mssqlDialect{}
}
