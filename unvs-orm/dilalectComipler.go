package orm

type DialectCompiler interface {
	resolve(aliasSource *map[string]string, caller *methodCall) (*resolverResult, error)
	getQuoteIdent() string
	getParam(index int) string
	driverName() string
	setCompiler(compiler *CompilerUtils)
	setJoinCompiler(compiler *JoinCompilerUtils)
	getCompiler() *CompilerUtils
}
