package internal

import "reflect"

var utils = &utilsPackage{}

func InitUtils(currentPackagePath string, entityTypeName string, mapType map[reflect.Type]string) {
	utils.currentPackagePath = currentPackagePath
	utils.entityTypeName = entityTypeName
	utils.mapType = mapType
}

var Utils = utils
