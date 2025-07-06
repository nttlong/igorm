package orm

type expression struct {
	dialect     DialectCompiler
	cmp         *CompilerUtils
	keywords    []string
	specialChar []byte
}
