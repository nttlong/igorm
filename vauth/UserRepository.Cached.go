package unvsauth

import "unvs-auth/models"

type UserRepositoryCached struct {
	db    UserRepository
	cache UserCache
}

func (r *UserRepositoryCached) FindByEmailOrUsername(id string) (*models.User, error) {
	if user, ok := r.cache.Get(id); ok {
		return user, nil
	}
	user, err := r.db.FindByEmailOrUsername(id)
	if err == nil {
		r.cache.Set(user)
	}
	return user, err
}

func (r *UserRepositoryCached) Create(u *models.User) error {
	err := r.db.Create(u)
	if err == nil {
		r.cache.Set(u)
	}
	return err
}

func (r *UserRepositoryCached) Delete(id string) error {
	err := r.db.Delete(id)
	if err == nil {
		r.cache.Delete(id)
	}
	return err
}
