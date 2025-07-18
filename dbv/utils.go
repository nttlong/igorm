package dbv

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
type initPlural struct {
	once sync.Once
	val  string
}

func (u *utilsReceiver) Plural(txt string) string {
	actual, _ := u.cachePlural.LoadOrStore(txt, &initPlural{})
	init := actual.(*initPlural)

	init.once.Do(func() {
		txt = strings.ToLower(txt)
		ret := pluralize.Plural(txt)
		init.val = ret
	})
	return init.val
}

type initToSnakeCase struct {
	once sync.Once
	val  string
}

func (u *utilsReceiver) ToSnakeCase(str string) string {
	actual, _ := u.cacheToSnakeCase.LoadOrStore(str, &initToSnakeCase{})
	init := actual.(*initToSnakeCase)

	init.once.Do(func() {
		init.val = u.toSnakeCase(str)
	})
	return init.val
}
func (u *utilsReceiver) toSnakeCase(str string) string {

	var result []rune
	for i, r := range str {
		if i > 0 && unicode.IsUpper(r) &&
			(unicode.IsLower(rune(str[i-1])) || (i+1 < len(str) && unicode.IsLower(rune(str[i+1])))) {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}
	ret := string(result)

	return ret
}

var utils = &utilsReceiver{
	cacheToSnakeCase: sync.Map{},
	cachePlural:      sync.Map{},
	cacheQueryable:   sync.Map{},
	cacheQuoteText:   sync.Map{},
	EXPR:             *newExprUtils(),
}
