package internal

import (
	"encoding/json"
	"fmt"
	"reflect"
)

/*
This struct is only used for the function buildRepositoryFromType of the "utilsPackage"
*/
type repositoryValueStruct struct {
	/*
		The buildRepositoryFromType function of utilsPackage
		will analyze the type of repo: During the analysis process,
		it will use reflect.New(type) to create a value for this field
	*/
	ValueOfRepo    reflect.Value
	PtrValueOfRepo reflect.Value
	/*
		During the analysis of the entity type, these are fields that have struct types,
		and those structs have fields declared with types like DbField[<type>], including in embedded struct
	*/
	EntityTypes []reflect.Type
}
type autoNumberKey struct {
	FieldName string
	fieldTag  *FieldTag
	KeyType   reflect.Type
}

// --------------------- Tag Metadata ---------------------
// FieldTag holds parsed metadata from struct field tags.

type fkInfo struct {
	FromTable string
	FromField []string
	ToTable   string
	ToField   []string
}

type schemaMap struct {
	table  map[string]bool
	unique map[string]bool
	index  map[string]bool
	fk     map[string]bool
}
type Model[T any] struct {
	Meta map[string]FieldTag
}

func ToJsonString(data interface{}) string {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return ""
	}
	return string(jsonData)
}
func PrintJson(data interface{}) {
	fmt.Println(ToJsonString(data))
}
