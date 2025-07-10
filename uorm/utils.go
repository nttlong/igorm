package uorm

import (
	"reflect"
	"sync"
	"unicode"

	pluralizeLib "github.com/gertd/go-pluralize"
)

var pluralize = pluralizeLib.NewClient()

type utilsReceiver struct {
	cacheToSnakeCase sync.Map
	cachePlural      sync.Map
	cacheQueryable   sync.Map
	cacheQuoteText   sync.Map
}

func (u *utilsReceiver) Plural(txt string) string {
	if v, ok := u.cachePlural.Load(txt); ok {
		return v.(string)
	}
	ret := pluralize.Plural(txt)
	u.cachePlural.Store(txt, ret)

	return ret
}
func (u *utilsReceiver) ToSnakeCase(str string) string {
	if v, ok := u.cacheToSnakeCase.Load(str); ok {
		return v.(string)
	}
	var result []rune
	for i, r := range str {
		if i > 0 && unicode.IsUpper(r) &&
			(unicode.IsLower(rune(str[i-1])) || (i+1 < len(str) && unicode.IsLower(rune(str[i+1])))) {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}
	ret := string(result)
	u.cacheToSnakeCase.Store(str, ret)
	return ret
}
func (u *utilsReceiver) QuoteText(dbType DB_TYPE, text string) string {
	key := dbType.String() + ":" + text
	if v, ok := u.cacheQuoteText.Load(key); ok {
		return v.(string)
	}
	ret := ""
	switch dbType {
	case DB_TYPE_MYSQL:
		ret = "`" + text + "`"
	case DB_TYPE_POSTGRES:
		ret = "\"" + text + "\""
	case DB_TYPE_SQLITE:
		ret = "\"" + text + "\""
	case DB_TYPE_MSSQL:
		ret = "[" + text + "]"
	default:
		ret = text
	}
	u.cacheQuoteText.Store(key, ret)
	return ret

}
func (u *utilsReceiver) createQueryableFormType(typ reflect.Type, tableName string, dbType DB_TYPE, isAlias bool) interface{} {
	key := typ.String() + "://" + tableName
	if v, ok := u.cacheQueryable.Load(key); ok {
		return v
	}
	obj := reflect.New(typ).Elem()
	strTableName := ""
	if isAlias {
		strTableName = u.QuoteText(dbType, tableName)
	} else {
		strTableName = u.QuoteText(dbType, u.Plural(u.ToSnakeCase(tableName)))
	}
	model := Model{
		entity: typ,
		table: &Table{
			name: strTableName,
		},
		dbType: dbType,
	}
	obj.FieldByName("Model").Set(reflect.ValueOf(model))
	for i := 0; i < typ.NumField(); i++ {
		if typ.Field(i).Type == reflect.TypeOf(Model{}) {
			continue
		}
		newField := Field{
			expr:  u.QuoteText(dbType, u.ToSnakeCase(typ.Field(i).Name)),
			name:  u.QuoteText(dbType, typ.Field(i).Name),
			table: model.table,
			args:  nil,
		}
		fieldVal := obj.FieldByName(typ.Field(i).Name)
		if fieldVal.CanSet() {
			flf := obj.FieldByName(typ.Field(i).Name)
			val := reflect.ValueOf(newField)
			if flf.Type().Kind() == reflect.Ptr {
				flf.Set(val.Addr())
			} else {
				flf.Set(val)
			}

		} else {
			panic("field " + typ.Field(i).Name + " is not settable")
		}

	}
	ret := obj.Interface()
	u.cacheQueryable.Store(key, ret)
	return ret
}

var utils = utilsReceiver{}
