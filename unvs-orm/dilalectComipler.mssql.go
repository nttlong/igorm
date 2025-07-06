package orm

import (
	"fmt"
	"strconv"
	"strings"
)

type mssqlDialect struct {
	compiler     *CompilerUtils
	joinCompiler *JoinCompilerUtils
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
func (d *mssqlDialect) setJoinCompiler(compiler *JoinCompilerUtils) {
	d.joinCompiler = compiler
}
func (d *mssqlDialect) getCompiler() *CompilerUtils {
	if d.compiler == nil {
		d.compiler = Compiler.Ctx(d)
	}
	return d.compiler
}
func (d *mssqlDialect) resolve(aliasSource *map[string]string, caller *methodCall) (*resolverResult, error) {
	if caller.method == "text" {
		return d.textFunc(aliasSource, caller)
	}
	strArgs := make([]string, 0)
	retArgs := make([]interface{}, 0)
	if caller.dbField != nil {
		field, err := d.compiler.Resolve(aliasSource, caller.dbField)
		if err != nil {
			return nil, err
		}
		strArgs = append(strArgs, field.Syntax)
		retArgs = append(retArgs, field.Args...)
	}
	for _, arg := range caller.args {
		if strArf, ok := arg.(string); ok {
			strArgs = append(strArgs, strArf)
			continue
		} else {
			rs, err := d.compiler.Resolve(aliasSource, arg)
			if err != nil {
				return nil, err
			}
			strArgs = append(strArgs, rs.Syntax)
			retArgs = append(retArgs, rs.Args...)
		}
	}

	return &resolverResult{
		Syntax: caller.method + "(" + strings.Join(strArgs, ", ") + ")",
		Args:   retArgs,
	}, nil
}
func (d *mssqlDialect) textFunc(aliasSource *map[string]string, caller *methodCall) (*resolverResult, error) {
	//CONVERT(NVARCHAR(50), 12345)
	if len(caller.args) != 1 {
		return nil, fmt.Errorf("text function only accept one argument")
	}
	arg := caller.args[0]
	txtArgs := ""
	args := make([]interface{}, 0)
	if strArf, ok := arg.(string); ok {
		txtArgs = strArf
	} else {
		rs, err := d.compiler.Resolve(aliasSource, arg)
		if err != nil {
			return nil, err
		}
		txtArgs = rs.Syntax
		args = append(args, rs.Args...)
	}
	return &resolverResult{
		Syntax: "CONVERT(NVARCHAR(50), " + txtArgs + ")",
		Args:   args,
	}, nil
}

func (d *mssqlDialect) driverName() string {
	return "mssql"
}

var MssqlDialect = mssqlDialect{}

func NewMssqlDialect() DialectCompiler {
	return &mssqlDialect{}
}
