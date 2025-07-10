package orm

func (sql *SqlCmdSelect) buildExprField(f *exprField) (string, []interface{}, error) {

	ret, err := sql.cmp.Resolve(sql.tables, sql.buildContext, f, false, true)

	if err != nil {
		return "", nil, err
	}
	return ret.Syntax, ret.Args, nil
}
