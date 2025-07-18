package unvsauth

import (
	"database/sql"
	"dbv"
	_ "dbv"
	"errors"

	models "unvs-auth/models"
)

type UserRepositorySQL struct {
	db *dbv.TenantDB
}

func NewUserRepositorySQL(db *dbv.TenantDB) *UserRepositorySQL {
	return &UserRepositorySQL{db: db}
}

func (r *UserRepositorySQL) FindByEmailOrUsername(identifier string) (*User, error) {

	dbv.SelectAll[User]()
	user := models.User{}

	r.db.First(&user, "email = ? OR username = ?", identifier, identifier)

	query := `SELECT id, user_id, email, username, hash_password, is_active
	          FROM users
			  WHERE email = ? OR username = ? LIMIT 1`
	row := r.db.QueryRow(query, identifier, identifier)

	var u User
	err := row.Scan(&u.ID, &u.UserId, &u.Email, &u.Username, &u.HashPassword, &u.IsActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepositorySQL) Create(u *User) error {
	query := `INSERT INTO users (user_id, email, username, hash_password, is_active)
	          VALUES (?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, u.UserId, u.Email, u.Username, u.HashPassword, u.IsActive)
	return err
}
