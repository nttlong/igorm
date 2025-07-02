package internal

import "fmt"

/*
This function will read information from @typ and create a structure similar to the one described in
"repositoryValueStruct" if no error occurs
*/
type buildRepositoryError struct {
	FieldName     string
	FieldTypeName string
	err           error
}

func (e buildRepositoryError) Error() string {
	return fmt.Sprintf("build repository from type %s error: %s", e.FieldTypeName, e.FieldName)
}
