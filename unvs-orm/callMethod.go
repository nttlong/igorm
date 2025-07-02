package orm

type methodCall struct {
	*dbField
	method string

	args []interface{}
}
