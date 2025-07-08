package orm

type DialectCompiler interface {
	resolve(tables *[]string, context *map[string]string, caller *methodCall, requireAlias bool) (*resolverResult, error)
	getQuoteIdent() string
	getParam(index int) string
	driverName() string
	setCompiler(compiler *CompilerUtils)
	setJoinCompiler(compiler *JoinCompilerUtils)
	getCompiler() *CompilerUtils
}
