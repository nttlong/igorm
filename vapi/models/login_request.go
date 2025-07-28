// models/login_request.go
package models

type LoginRequest struct {
	Username string `json:"username" form:"username" example:"admin"`
	Password string `json:"password" form:"password" example:"123456"`
}
