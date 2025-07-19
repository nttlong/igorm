package vauth

import (
	"vdb"
	_ "vdb"

	models "vauth/models"
)

type UserRepositorySQL struct {
	db *vdb.TenantDB
}

func NewUserRepositorySQL(db *vdb.TenantDB) *UserRepositorySQL {
	return &UserRepositorySQL{db: db}
}

func (r *UserRepositorySQL) FindByEmailOrUsername(identifier string) (*models.User, error) {

	user := models.User{}

	err := r.db.First(&user, "email = ? OR username = ?", identifier, identifier)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepositorySQL) Create(u *models.User) error {
	err := r.db.Insert(u)
	return err
}
func (r *UserRepositorySQL) Delete(id string) error {
	_, err := r.db.Delete(&models.User{}, "userId = ?", id)
	return err
}
