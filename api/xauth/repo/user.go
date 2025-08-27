package repo

import (
	"errors"
	"sync"
	"time"
	"vdb"
	"xauth/models"
)

type UserRepo interface {
	GetUser() (*models.User, error)
	CreateDefautUser(hashPassword string) error
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
}
type UserRepoSql struct {
	db *vdb.TenantDB
}

func NewUserRepoSql(db *vdb.TenantDB) UserRepo {
	return &UserRepoSql{
		db: db,
	}
}
func (u *UserRepoSql) GetUser() (*models.User, error) {

	panic("Not imeplemented")
}

var createDefautUserOnce sync.Once

func (u *UserRepoSql) CreateDefautUser(hasnPassword string) error {
	var err error
	createDefautUserOnce.Do(func() {
		err = u.CreateUser(&models.User{
			Username:  "admin",
			Password:  hasnPassword,
			Active:    true,
			Email:     nil,
			Phone:     nil,
			CreatedOn: time.Now().UTC(),
		})
		var dbErr *vdb.DialectError
		if err != nil && errors.As(err, &dbErr) {

			if dbErr.ErrorType == vdb.DIALECT_DB_ERROR_TYPE_DUPLICATE {
				// User is existing accept, no error
				err = nil
			}
		}

	})

	return err

}
func (u *UserRepoSql) CreateUser(user *models.User) error {
	user.Active = true
	user.CreatedOn = time.Now().UTC()
	err := u.db.Create(user)
	return err

}
func (u *UserRepoSql) UpdateUser(user *models.User) error {
	ret := u.db.Update(user)
	return ret.Error

}
