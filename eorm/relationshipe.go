package eorm

import (
	"eorm/migrate"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Relationship struct {
}

func (r *Relationship) BelongsTo(model interface{}, foreignKey string) *Relationship {
	return nil
}
func (m *Model[T]) getPrimaryKeys() []migrate.ColumnDef {
	info := ModelRegistry.GetModelByType(reflect.TypeFor[T]())
	for _, pk := range info.GetPrimaryConstraints() {
		return pk
	}
	return nil
}
func (m *Model[T]) AddForeignKey(FkEntity interface{}, foreignKey, keys string) *Model[T] {
	ks := strings.Split(keys, ",")
	fks := strings.Split(foreignKey, ",")
	ownerType := reflect.TypeFor[T]()
	FkEntityType := reflect.TypeOf(FkEntity)
	if FkEntityType.Kind() == reflect.Ptr {
		FkEntityType = FkEntityType.Elem()
	}
	ownerInfo := ModelRegistry.GetModelByType(ownerType)
	fkInfo := ModelRegistry.GetModelByType(FkEntityType)
	if FkEntityType.Kind() == reflect.Ptr {
		FkEntityType = FkEntityType.Elem()
	}

	if len(ks) != len(fks) {
		panic(fmt.Sprintf("len of key and foreign key not match: %s(%s)!= %s(%s)", ownerInfo.GetTableName(), keys, fkInfo.GetTableName(), foreignKey))
	}
	ownerMap := map[string]migrate.ColumnDef{}
	for _, col := range ownerInfo.GetColumns() {
		ownerMap[strings.ToLower(col.Field.Name)] = col
	}
	fkMap := map[string]migrate.ColumnDef{}

	for _, col := range fkInfo.GetColumns() {
		fkMap[strings.ToLower(col.Field.Name)] = col

	}
	pkCols := []string{}
	fkColsName := []string{}
	for i, key := range ks {
		keyCol := ownerMap[strings.ToLower(key)]
		fkCol := fkMap[strings.ToLower(fks[i])]
		if keyCol.Type != fkCol.Type {
			errText := fmt.Sprintf("foreign key column not match with primary key of %s.%s and %s.%s ", ownerType.String(), keyCol.Name, FkEntityType.String(), fkCol.Name)
			panic(errors.New(errText))
		}
		pkCols = append(pkCols, keyCol.Name)
		fkColsName = append(fkColsName, fkCol.Name)

	}

	ForeignKeyRegistry.Register(&foreignKeyInfo{
		fromTable: ownerInfo.GetTableName(),
		fromCols:  fkColsName,
		toTable:   fkInfo.GetTableName(),
		toCols:    fkColsName,
	})

	return m

}

type foreignKeyInfo struct {
	fromTable string
	fromCols  []string
	toTable   string
	toCols    []string
}
type foreignKeyRegistry struct {
	fkMap map[string]*foreignKeyInfo
}

func (r *foreignKeyRegistry) Register(fk *foreignKeyInfo) {
	key := fmt.Sprintf("FK_%s__%s_____%s__%s", fk.fromTable, strings.Join(fk.fromCols, "_____"), fk.toTable, strings.Join(fk.toCols, "_____"))
	r.fkMap[key] = fk

}

var ForeignKeyRegistry = foreignKeyRegistry{
	fkMap: map[string]*foreignKeyInfo{},
}
