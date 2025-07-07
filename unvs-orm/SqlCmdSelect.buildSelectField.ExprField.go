package orm

func (sql *SqlCmdSelect) buildExprField(f *exprField) (string, []interface{}, error) {

	ret, err := sql.cmp.Resolve(sql.buildContext, f)

	if err != nil {
		return "", nil, err
	}
	return ret.Syntax, ret.Args, nil
}
