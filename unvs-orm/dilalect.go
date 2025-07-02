package orm

type Dialect interface {
	resolve(caller *methodCall) (*resolverResult, error)
	getQuoteIdent() string
	getParam(index int) string
	driverName() string
	setCompiler(compiler *CompilerUtils)
	getCompiler() *CompilerUtils
}
