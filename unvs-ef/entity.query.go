package unvsef

func (t *TenantDb) Query() *Query {
	ret := NewQuery()
	ret.dialect = t.Dialect
	return ret

}
func (t *Query) SQLCommand() (string, []interface{}) {
	return t.ToSQL(t.dialect)
}
