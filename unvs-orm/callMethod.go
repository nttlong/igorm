package orm

type methodCall struct {
	isFromExr bool
	method    string

	args []interface{}
}
