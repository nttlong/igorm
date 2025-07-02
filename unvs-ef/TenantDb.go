package unvsef

import (
	"database/sql"
	"fmt"
	"reflect"
)

/*
This struct is definition of tenant db
*/
type TenantDb struct {
	DB                   sql.DB
	Dialect              Dialect
	DBType               DBType
	DBTypeName           string
	SqlMigrate           []string
	DbName               string
	Relationships        []*RelationshipRegister
	buildRepositoryError error
	Err                  error
}
type MigrateError struct {
	Errs          []error
	SqlCauseError []string
}

func (e *MigrateError) Error() string {
	ret := ""
	for i, sql := range e.SqlCauseError {
		ret += fmt.Sprintf("Error while executing sql: %s\n Cause: %s\n", sql, e.Errs[i].Error())
	}

	return ret
}

func (t *TenantDb) DoMigrate() error {
	retErr := &MigrateError{
		Errs:          []error{},
		SqlCauseError: []string{},
	}

	for _, sql := range t.SqlMigrate {
		_, err := t.DB.Exec(sql)
		if err != nil {
			retErr.Errs = append(retErr.Errs, err)
			retErr.SqlCauseError = append(retErr.SqlCauseError, sql)
			continue
		}
	}
	if len(retErr.Errs) > 0 {
		return retErr
	}
	return nil
}

type RelationshipRegister struct {
	owner      *TenantDb
	err        error
	fromTable  string
	fromFields []string
	toTable    string
	toFields   []string
}

func (r *RelationshipRegister) From(fields ...interface{}) *RelationshipRegister {
	for _, field := range fields {
		valField := reflect.ValueOf(field)
		tableNameField := valField.FieldByName("TableName")
		tableName := ""
		fieldName := ""
		if tableNameField.IsValid() {
			tableName = tableNameField.String()
		} else {
			r.err = fmt.Errorf("TableName not found in field at RelationshipRegister.From '%s'", field)
			r.owner.buildRepositoryError = r.err
			return r
		}
		if r.fromTable == "" {
			r.fromTable = tableName
		}

		fieldNameField := valField.FieldByName("ColName")

		if fieldNameField.IsValid() {
			fieldName = fieldNameField.String()
		} else {
			r.err = fmt.Errorf("TableName not found in field in field at RelationshipRegister.From '%s'", field)
			return r
		}
		if r.fromTable != "" && r.fromTable != tableName {

			r.owner.buildRepositoryError = fmt.Errorf("%s must belong to the same table name '%s' but got '%s'", fieldName, r.fromTable, tableName)
			return r
		}

		r.fromFields = append(r.fromFields, fieldName)
		r.fromTable = tableName

	}
	return r

	// panic("not implemented")
}
func (r *RelationshipRegister) To(fields ...interface{}) error {
	if r.err != nil {
		return r.err
	}
	for _, field := range fields {
		valField := reflect.ValueOf(field)
		tableNameField := valField.FieldByName("TableName")
		tableName := ""
		fieldName := ""
		if tableNameField.IsValid() {
			tableName = tableNameField.String()
		} else {
			r.owner.buildRepositoryError = fmt.Errorf("TableName not found in field at RelationshipRegister.To '%s'", field)
			return r.owner.buildRepositoryError
		}
		fieldNameField := valField.FieldByName("ColName")
		if fieldNameField.IsValid() {
			fieldName = fieldNameField.String()
		} else {
			r.owner.buildRepositoryError = fmt.Errorf("TableName not found in field at RelationshipRegister.To '%s'", field)
			return r.owner.buildRepositoryError
		}
		if r.toTable == "" {
			r.toTable = tableName
		}
		if r.toTable != tableName {
			r.err = fmt.Errorf("%s must belong to the same table name '%s' but got '%s'", fieldName, r.fromTable, tableName)
			return r.err
		}

	}
	// check the balance of num of field in from and to
	if len(r.fromFields) != len(r.toFields) {
		r.err = fmt.Errorf("number of fields were not balance, from %d to %d", len(r.fromFields), len(r.toFields))
		return r.err
	}

	return nil

}

/*
While create repository looks like

	type Repository struct {
		*TenantDb //<-- Should I user pointer here
		Articles  *Article
		Comments  *Comment
	}
		func (r *Repository) Init() { // the relationship set up herr
				r.NewRelationship().From(r.Articles.Id).To(r.Comments.ArticleId, r.Comments.Id)

		}
*/
func (t *TenantDb) NewRelationship() *RelationshipRegister {
	ret := &RelationshipRegister{
		owner:      t,
		fromFields: []string{},
		toFields:   []string{},
		fromTable:  "",
		toTable:    "",
	}
	t.Relationships = append(t.Relationships, ret)
	return ret

}
func (t *TenantDb) From(table interface{}) *Query {
	return &Query{
		table:    table,
		tenantDb: t,
	}
}
func (t *TenantDb) InsertInto(entity any) *InsertQuery {
	typ := reflect.TypeOf(entity)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	tableName := utils.TableNameFromStruct(typ)
	return &InsertQuery{
		tableName:  tableName,
		entityType: typ,
		dialect:    t.Dialect,
		tenantDb:   t,
	}
}
