package migrate

import (
	"fmt"
	"strings"
)

type ForeignKeyInfo struct {
	FromTable string
	FromCols  []string
	ToTable   string
	ToCols    []string
}
type foreignKeyRegistry struct {
	fkMap map[string]*ForeignKeyInfo
}

func (r *foreignKeyRegistry) Register(fk *ForeignKeyInfo) {
	key := fmt.Sprintf("FK_%s__%s_____%s__%s", fk.FromTable, strings.Join(fk.FromCols, "_____"), fk.ToTable, strings.Join(fk.ToCols, "_____"))
	r.fkMap[key] = fk

}

var ForeignKeyRegistry = foreignKeyRegistry{
	fkMap: map[string]*ForeignKeyInfo{},
}
