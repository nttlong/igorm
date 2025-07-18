package unvsauth

import "errors"

type User struct {
	ID           int
	UserId       string
	Email        string
	Username     string
	HashPassword string
	IsActive     bool
}

type UserRepository struct {
	users map[string]*User // giả lập
}

func NewUserRepository() *UserRepository {
	return &UserRepository{users: map[string]*User{}}
}

func (r *UserRepository) FindByEmailOrUsername(id string) (*User, error) {
	for _, u := range r.users {
		if u.Email == id || u.Username == id {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *UserRepository) Create(u *User) error {
	if _, err := r.FindByEmailOrUsername(u.Email); err == nil {
		return errors.New("user already exists")
	}
	r.users[u.UserId] = u
	return nil
}
