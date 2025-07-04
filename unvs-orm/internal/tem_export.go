package internal

import "reflect"

func BuildRepositoryFromType(typ reflect.Type) (*repositoryValueStruct, error) {
	return utils.buildRepositoryFromType(typ, false)
}
