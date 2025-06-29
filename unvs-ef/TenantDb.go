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
	DB                   *sql.DB
	Dialect              Dialect
	DBType               DBType
	DBTypeName           string
	SqlMigrate           []string
	DbName               string
	Relationships        []*RelationshipRegister
	buildRepositoryError error
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
	owner *TenantDb
	err   error
	from  map[string][]string
	to    map[string][]string
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
		fieldNameField := valField.FieldByName("ColName")
		if fieldNameField.IsValid() {
			fieldName = fieldNameField.String()
		} else {
			r.err = fmt.Errorf("TableName not found in field in field at RelationshipRegister.From '%s'", field)
			return r
		}
		if _, ok := r.from[tableName]; !ok {
			r.from[tableName] = []string{}
		}
		r.from[tableName] = append(r.from[tableName], fieldName)

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
		if _, ok := r.to[tableName]; !ok {
			r.to[tableName] = []string{}
		}
		r.to[tableName] = append(r.from[tableName], fieldName)

	}
	// check the balance of num of field in from and to
	for tableName, fields := range r.from {
		if _, ok := r.to[tableName]; !ok {
			r.owner.buildRepositoryError = fmt.Errorf("number of fields were not balance, from %d to %d", len(fields), 0)
			return r.owner.buildRepositoryError

		}
		if len(fields) != len(r.to[tableName]) {
			r.owner.buildRepositoryError = fmt.Errorf("number of fields were not balance, from %d to %d", len(fields), len(r.to[tableName]))
			return r.owner.buildRepositoryError

		}
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
		owner: t,
		from:  make(map[string][]string),
		to:    make(map[string][]string),
	}
	t.Relationships = append(t.Relationships, ret)
	return ret

}
