package eorm

import (
	"strings"
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
	EXPR             exprUtils
}

func (u *utilsReceiver) Plural(txt string) string {
	if v, ok := u.cachePlural.Load(txt); ok {
		return v.(string)
	}
	txt = strings.ToLower(txt)
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

var utils = &utilsReceiver{
	cacheToSnakeCase: sync.Map{},
	cachePlural:      sync.Map{},
	cacheQueryable:   sync.Map{},
	cacheQuoteText:   sync.Map{},
	EXPR:             *newExprUtils(),
}
