package wx

import (
	"wx/errors"
)

type errorFactoty struct {
}

var Errors errorFactoty

func (err *errorFactoty) RequireErr(field ...string) error {
	return &errors.RequireError{
		Fields:  field,
		Message: "required",
	}
}
func init() {
	Errors = errorFactoty{}

}
