package orm

import "strings"

func createDBField(fullPath string) *dbField {
	return &dbField{
		Name:  strings.Split(fullPath, ".")[1],
		Table: strings.Split(fullPath, ".")[0],
	}
}

func CreateDateTimeField(fullPath string) DateTimeField {
	return DateTimeField{
		underField: createDBField(fullPath),
	}

}
func CreateTextField(fullPath string) TextField {
	return TextField{
		underField: createDBField(fullPath),
	}
}
func CreateNumberField[T Number](fullPath string) NumberField[T] {
	return NumberField[T]{
		underField: createDBField(fullPath),
	}
}
func CreateBoolField(fullPath string) BoolField {
	return BoolField{
		underField: createDBField(fullPath),
	}
}
