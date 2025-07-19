package vdb

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"vdb/migrate"
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

type CascadeOption migrate.CascadeOption

func (m *Model[T]) AddForeignKey(foreignKey string, FkEntity interface{}, keys string, cascadeOption *CascadeOption) *Model[T] {
	if cascadeOption == nil {
		cascadeOption = &CascadeOption{
			OnDelete: true,
			OnUpdate: true,
		}
	}

	ks := strings.Split(keys, ",")
	fks := strings.Split(foreignKey, ",")
	FkEntityType := reflect.TypeFor[T]()
	ownerType := reflect.TypeOf(FkEntity)
	if ownerType.Kind() == reflect.Ptr {
		ownerType = ownerType.Elem()
	}
	ModelRegistry.RegisterType(ownerType)
	ModelRegistry.RegisterType(FkEntityType)

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

	migrate.ForeignKeyRegistry.Register(&migrate.ForeignKeyInfo{
		FromTable: fkInfo.GetTableName(),
		FromCols:  fkColsName,
		ToTable:   ownerInfo.GetTableName(),
		ToCols:    pkCols,
	})

	return m

}
