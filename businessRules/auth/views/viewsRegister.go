package views

import (
	_ "dbmodels/auth"
	dbmodels "dbmodels/auth"
	"dbx"
	"dynacall"
)

// This service register a new view a view is a feature com
func (view *ViewService) Register(data struct {
	ViewId      string
	Name        string
	Description string
}) {
	view.TenantDb.Insert(&dbmodels.View{
		ViewId:      data.ViewId,
		Name:        data.Name,
		Description: dbx.FullTextSearchColumn(data.Description),
	})

}
func init() {
	dynacall.RegisterCaller(ViewService{
		Caller: dynacall.Caller{
			Path: "auth",
		},
	})
}
