package orm

import "sync"

type Join struct {
}

type JoinCompilerUtils struct {
	dialect DialectCompiler
}

var cacheJoinCompilerCtx sync.Map

func (c *JoinCompilerUtils) Ctx(dialect DialectCompiler) *JoinCompilerUtils {
	if dialect == nil {
		panic("dialect is nil")
	}

	key := dialect.driverName()

	if v, ok := cacheJoinCompilerCtx.Load(key); ok {
		return v.(*JoinCompilerUtils)
	}
	ret := &JoinCompilerUtils{dialect: dialect}
	dialect.setJoinCompiler(ret)
	cacheJoinCompilerCtx.Store(key, ret)
	return ret
}
func (c *JoinCompilerUtils) Resolve(expr interface{}) (*resolverResult, error) {
	panic("not implemented in JoinCompilerUtils.Resolve, row 30 file unvs-orm/SqlBuildrJoin.go")
}

var JoinCompiler = JoinCompilerUtils{}
