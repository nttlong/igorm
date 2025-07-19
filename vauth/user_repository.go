package unvsauth

import (
	"errors"
	"unvs-auth/models"
)

type UserRepository interface {
	FindByEmailOrUsername(identifier string) (*models.User, error)
	Create(user *models.User) error
	Delete(userID string) error // nếu bạn có logic xóa, có thể bỏ nếu chưa dùng
}
type MockUserRepository struct {
	users map[string]*models.User // giả lập
}

func NewUserRepository() *MockUserRepository {
	return &MockUserRepository{users: map[string]*models.User{}}
}

func (r *MockUserRepository) FindByEmailOrUsername(id string) (*models.User, error) {
	for _, u := range r.users {
		if u.Email == id || u.Username == id {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *MockUserRepository) Create(u *models.User) error {
	if _, err := r.FindByEmailOrUsername(u.Email); err == nil {
		return errors.New("user already exists")
	}
	r.users[u.UserId] = u
	return nil
}
