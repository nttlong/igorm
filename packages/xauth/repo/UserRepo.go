package repo

import (
	"fmt"
	"wx"
)

type UserRepo interface {
}
type UserRepoSQL struct {
}

func (userRepo *UserRepoSQL) New(ctx *wx.HttpContext, dbContext *wx.HttpService[DbContext]) (UserRepo, error) {
	fmt.Println(dbContext)
	return &UserRepoSQL{}, nil

}
