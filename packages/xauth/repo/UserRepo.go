package repo

type UserRepo interface {
}
type UserRepoSQL struct {
	db *DbContext
}

func NewUserRepoSQL(db *DbContext) *UserRepoSQL {
	return &UserRepoSQL{
		db: db,
	}
}
