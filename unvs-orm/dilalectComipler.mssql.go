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

	return d.compiler
}
func (d *mssqlDialect) resolve(tables *[]string, context *map[string]string, caller *methodCall, extractAlias, applyContext bool) (*resolverResult, error) {

	if caller.isFromExr {
		strArgs := []string{}
		for _, arg := range caller.args {
			if strArg, ok := arg.(string); ok {
				strArgs = append(strArgs, strArg)
			} else {
				return nil, fmt.Errorf("syntax error: unsupported argument type: %T", arg)
			}
		}
		return &resolverResult{
			Syntax: fmt.Sprintf("%s (%s)", caller.method, strings.Join(strArgs, ", ")),
		}, nil
	}
	methodName := strings.ToLower(caller.method)
	if methodName == "text" {
		return d.textFunc(tables, context, caller, extractAlias, applyContext)
	}
	strArgs := make([]string, 0)
	retArgs := make([]interface{}, 0)

	for _, arg := range caller.args {
		rs, err := d.compiler.Resolve(tables, context, arg, extractAlias, applyContext)
		if err != nil {
			return nil, err
		}
		strArgs = append(strArgs, rs.Syntax)
		retArgs = append(retArgs, rs.Args...)
	}
	if methodName == "format" {
		return &resolverResult{
			Syntax: caller.method + "(" + strArgs[0] + ",?" + ")",
			Args:   retArgs,
		}, nil
	}
	return &resolverResult{
		Syntax: caller.method + "(" + strings.Join(strArgs, ", ") + ")",
		Args:   retArgs,
	}, nil
}
func (d *mssqlDialect) textFunc(tables *[]string, context *map[string]string, caller *methodCall, extractAlias, applyContext bool) (*resolverResult, error) {
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
		rs, err := d.compiler.Resolve(tables, context, arg, extractAlias, applyContext)
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

var MssqlDialect = mssqlDialect{
	compiler:     &CompilerUtils{},
	joinCompiler: &JoinCompilerUtils{},
}
var MssqlCompiler = &CompilerUtils{}
var MssqlJoinCompilerUtils = &JoinCompilerUtils{}

func init() {
	MssqlDialect.compiler = MssqlCompiler
	MssqlDialect.joinCompiler = MssqlJoinCompilerUtils
	MssqlCompiler.dialect = &MssqlDialect
	MssqlJoinCompilerUtils.dialect = &MssqlDialect
}
