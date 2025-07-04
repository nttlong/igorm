package internal

import (
	"database/sql"
	"fmt"
)

/*
This struct is definition of tenant db
*/
type TenantDb struct {
	*sql.DB
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
