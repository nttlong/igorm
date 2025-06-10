package users

import (
	"dbx"
	userModel "unvs/internal/model/auth"
	"unvs/views"
)

type UserFilter struct {
	Filter string        `json:"filter"`
	Params []interface{} `json:"params"`
	Sort   string        `json:"sort"`
	Page   int           `json:"page"`
	Size   int           `json:"size"`
}
type UserItem struct {
	userModel.User
	Index int `json:"index"`
}

// this is business logic for creating user
func (v *User) List(Filter *UserFilter) (*Response, error) {
	qr := dbx.Pager[userModel.User](&v.DbTenant, v.Context)
	if Filter.Filter != "" {
		qr.Where(Filter.Filter, Filter.Params...)
	}
	if Filter.Size == 0 {
		Filter.Size = 100
	}
	qr.Size(Filter.Size)
	qr.Page(Filter.Page)
	if Filter.Sort != "" {
		qr.Sort(Filter.Sort)
	}

	// ret, err := qr.Items()
	// if err != nil {
	// 	return &Response{Error: err}, nil
	// } else {
	// 	return &Response{Data: ret}, nil
	// }
	r, e := qr.Query()
	if e != nil {
		return &Response{Error: e}, nil
	} else {
		return &Response{Data: r}, nil
	}

}
func init() {
	views.AddView(&User{
		BaseView: views.BaseView{
			ViewPath: "auth/users",
			IsAuth:   true,
		},
	})
}
